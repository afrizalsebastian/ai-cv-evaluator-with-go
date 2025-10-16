package service_consumer

import (
	"context"
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/domain/models"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/domain/models/dao"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/domain/repository"
	chromaclient "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/chroma-client"
	geminiclient "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/gemini-client"
	ingestdocument "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/ingest-document"
)

type ICvEvaluatorConsumerService interface {
	RunningJob(ctx context.Context, jobId string) error
}

type cvEvaluatorConsumerService struct {
	gemini      geminiclient.IGeminiClient
	chroma      chromaclient.IChromaClient
	ingest      ingestdocument.IIngestFile
	cvEvaluator repository.ICvEvaluatorJobRepository
}

func NewCvEvaluatorConsumerService(
	gemini geminiclient.IGeminiClient,
	chroma chromaclient.IChromaClient,
	ingest ingestdocument.IIngestFile,
	cvEvaluator repository.ICvEvaluatorJobRepository,
) ICvEvaluatorConsumerService {
	return &cvEvaluatorConsumerService{
		gemini:      gemini,
		chroma:      chroma,
		ingest:      ingest,
		cvEvaluator: cvEvaluator,
	}
}

func (c *cvEvaluatorConsumerService) RunningJob(ctx context.Context, jobId string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	// Todo: Get JOB ITEM FROM DB
	job, err := c.cvEvaluator.GetByJobId(ctx, jobId)
	if err != nil {
		log.Println("failed to get job item")
		return err
	}

	// update to processing
	job.Status = models.StatusProcessing
	c.cvEvaluator.UpdateJobByJobId(ctx, jobId, job)

	// extract text from file
	extractedCv, err := c.ingest.ExtractTextFromPdf(path.Join("uploaded-file", job.FileId, "cv_file.pdf"))
	if err != nil {
		c.jobFailToProcess(ctx, job, err)
		return err
	}

	extractedReport, err := c.ingest.ExtractTextFromPdf(path.Join("uploaded-file", job.FileId, "report_file.pdf"))
	if err != nil {
		c.jobFailToProcess(ctx, job, err)
		return err
	}

	// Evaluate CV
	jobDescription, err := c.chroma.Query(ctx, "job_description", job.JobTitle+" "+"job description", 5)
	if err != nil {
		c.jobFailToProcess(ctx, job, err)
		return err
	}
	cvRubric, err := c.chroma.Query(ctx, "cv_rubric", job.JobTitle+" "+"cv rubric", 5)
	if err != nil {
		c.jobFailToProcess(ctx, job, err)
		return err
	}

	cvEvaluatePrompt := c.buildCvEvaluatorPrompt(job.JobTitle, extractedCv, jobDescription, cvRubric)
	cvGeminiResp, err := c.gemini.GenerateContent(ctx, job.JobTitle, cvEvaluatePrompt)
	if err != nil {
		c.jobFailToProcess(ctx, job, err)
		return err
	}
	cvResult := strings.Split(cvGeminiResp, "\n---\n")
	if len(cvResult) < 2 {
		err = fmt.Errorf("invalid response from gemini")
		c.jobFailToProcess(ctx, job, err)
		return err
	}
	job.CvMatchRate = cvResult[0]
	job.CvFeedback = cvResult[1]
	fmt.Println("job with id " + job.JobId + " have done processed cv")

	// Evaluate Report
	caseStudyBrief, err := c.chroma.Query(ctx, "case_study_brief", job.JobTitle+" "+"case study brief", 5)
	if err != nil {
		c.jobFailToProcess(ctx, job, err)
		return err
	}

	reportRubric, err := c.chroma.Query(ctx, "project_report_rubric", job.JobTitle+" "+"project report rubric", 5)
	if err != nil {
		c.jobFailToProcess(ctx, job, err)
		return err
	}

	reportEvaluatePrompt := c.buildReportEvaluatorPrompt(job.JobTitle, extractedReport, caseStudyBrief, reportRubric)
	reportGeminiResp, err := c.gemini.GenerateContent(ctx, job.JobTitle, reportEvaluatePrompt)
	if err != nil {
		c.jobFailToProcess(ctx, job, err)
		return err
	}
	reportResult := strings.Split(reportGeminiResp, "\n---\n")
	if len(reportResult) < 2 {
		err = fmt.Errorf("invalid response from gemini")
		c.jobFailToProcess(ctx, job, err)
		return err
	}
	job.ProjectScore = reportResult[0]
	job.ProjectFeedback = reportResult[1]
	fmt.Println("job with id " + job.JobId + " have done processed report")

	// final
	finalPrompt := c.buildFinalPrompt(job.CvMatchRate, job.CvFeedback, job.ProjectScore, job.ProjectFeedback)
	overall, err := c.gemini.GenerateContent(ctx, job.JobTitle, finalPrompt)
	if err != nil {
		c.jobFailToProcess(ctx, job, err)
		return err
	}
	job.OverallSummary = overall
	job.Status = models.StatusCompleted
	c.cvEvaluator.UpdateJobByJobId(ctx, jobId, job)

	return nil
}

func (w *cvEvaluatorConsumerService) jobFailToProcess(ctx context.Context, job *dao.CvEvaluatorJob, err error) {
	fmt.Printf("job with id %s failed to process: %s\n", job.JobId, err.Error())
	job.Status = models.StatusFailed
	_ = w.cvEvaluator.UpdateJobByJobId(ctx, job.JobId, job)
}

func (w *cvEvaluatorConsumerService) buildCvEvaluatorPrompt(jobTitle, extractedCv string, jobDescription, cvRubric []models.ChromaSearchResult) string {
	prompt := "Evaluate this CV for role: " + jobTitle + "\n"
	prompt += "Job Description: \n"
	for _, desc := range jobDescription {
		prompt += desc.Text
		prompt += "\n"
	}
	prompt += "\n-----\n"
	prompt += "CV Rubric: \n"
	for _, rubric := range cvRubric {
		prompt += rubric.Text
		prompt += "\n"
	}
	prompt += "\n----\n"
	prompt += "With Candidate CV: \n" + extractedCv
	prompt += "\n-----\n"
	prompt += "Return as:\n<0.0-1.0 match rate>\n---\n<brief feedback with 2-3 sentences>\n"
	return prompt
}

func (w *cvEvaluatorConsumerService) buildReportEvaluatorPrompt(jobTitle, extractedReport string, studyCase, reportRubric []models.ChromaSearchResult) string {
	prompt := "Evaluate this Project report for role: " + jobTitle + "\n"
	prompt += "Role study case: \n"
	for _, desc := range studyCase {
		prompt += desc.Text
		prompt += "\n"
	}
	prompt += "\n-----\n"
	prompt += "Projec Report Rubric: \n"
	for _, rubric := range reportRubric {
		prompt += rubric.Text
		prompt += "\n"
	}
	prompt += "\n----\n"
	prompt += "With Candidate Project Report: \n" + extractedReport
	prompt += "\n-----\n"
	prompt += "Return as:\n<1.0-5.0 project score>\n---\n<brief feedback with 2-3 sentences>\n"
	return prompt
}

func (w *cvEvaluatorConsumerService) buildFinalPrompt(cvRate, cvFeedback, projScore, projFeedback string) string {
	prompt := "Give 3-5 sentences summary based on:\n"
	prompt += "CV match rate: " + cvRate + "\n"
	prompt += "CV feedback: " + cvFeedback + "\n"
	prompt += "Project score: " + projScore + "\n"
	prompt += "Project feedback: " + projFeedback + "\n"
	prompt += "\nOutput concise summary (strengths, gaps, recommendations, advice, and other positive thing to improvement)."
	prompt += "Return as:\n<3-5 sentences for summary>"
	return prompt
}
