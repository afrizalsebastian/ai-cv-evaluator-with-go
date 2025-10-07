package worker

import (
	"context"
	"fmt"
	"log"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/services"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/models"
	chromaclient "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/chroma-client"
	geminiclient "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/gemini-client"
)

type ICvEvaluatorWorker interface {
	Start(ctx context.Context)
	Enqueue(job *models.JobItem) error
	Stop(ctx context.Context)
}

type cvEvaluatorWorker struct {
	fileStore services.ILocalFileStore
	gemini    geminiclient.IGeminiClient
	chroma    chromaclient.IChromaClient
	ingest    services.IIngestFile

	queue chan *models.JobItem
	mu    sync.RWMutex
	wg    sync.WaitGroup

	maxWorkers int
}

func NewCvEvaluatorWorker(
	fileStore services.ILocalFileStore,
	gemini geminiclient.IGeminiClient,
	chroma chromaclient.IChromaClient,
	ingest services.IIngestFile,
	maxWorkers int,
) ICvEvaluatorWorker {
	if maxWorkers <= 0 {
		maxWorkers = 5
	}

	return &cvEvaluatorWorker{
		fileStore:  fileStore,
		gemini:     gemini,
		chroma:     chroma,
		ingest:     ingest,
		maxWorkers: maxWorkers,
		queue:      make(chan *models.JobItem, 100),
	}
}

func (w *cvEvaluatorWorker) Start(ctx context.Context) {
	log.Printf("Worker start in background")
	for i := 0; i < w.maxWorkers; i++ {
		w.wg.Add(1)
		go w.worker(ctx)
	}
}

func (w *cvEvaluatorWorker) Enqueue(job *models.JobItem) error {
	select {
	case w.queue <- job:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("queue full, timeout enqueueing job")
	}
}

func (w *cvEvaluatorWorker) Stop(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	close(w.queue)
	w.wg.Wait()
}

func (w *cvEvaluatorWorker) worker(ctx context.Context) {
	defer w.wg.Done()
	for {
		select {
		case job, ok := <-w.queue:
			if !ok {
				return
			}
			w.processJob(ctx, job)
		case <-ctx.Done():
			return
		}
	}
}

func (w *cvEvaluatorWorker) processJob(ctx context.Context, job *models.JobItem) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	job.Status = models.StatusProcessing

	// extract text from file
	extractedCv, err := w.ingest.ExtractTextFromPdf(path.Join("uploaded-file", job.FileId, "cv_file.pdf"))
	if err != nil {
		w.jobFailToProcess(job, err)
		return
	}

	_, err = w.ingest.ExtractTextFromPdf(path.Join("uploaded-file", job.FileId, "report_file.pdf"))
	if err != nil {
		w.jobFailToProcess(job, err)
		return
	}

	jobDescription, err := w.chroma.Query(ctx, "job_description", job.JobTitle+" "+"job description", 5)
	if err != nil {
		w.jobFailToProcess(job, err)
		return
	}
	cvRubric, err := w.chroma.Query(ctx, "cv_rubric", job.JobTitle+" "+"cv rubric", 5)
	if err != nil {
		w.jobFailToProcess(job, err)
		return
	}

	cvEvaluatePrompt := w.buildCvEvaluatorPrompt(job.JobTitle, extractedCv, jobDescription, cvRubric)
	cvGeminiResp, err := w.gemini.GenerateContent(ctx, job.JobTitle, cvEvaluatePrompt)
	if err != nil {
		w.jobFailToProcess(job, err)
		return
	}
	cvResult := strings.Split(cvGeminiResp, "\n---\n")
	fmt.Println(cvResult)

	// caseStudyBrief, err := w.chroma.Query(ctx, "case_study_brief", job.JobTitle+" "+"case study brief", 5)
	// if err != nil {
	// 	w.jobFailToProcess(job, err)
	// 	return
	// }

	// reportRubric, err := w.chroma.Query(ctx, "project_report_rubric", job.JobTitle+" "+"project report rubric", 5)
	// if err != nil {
	// 	w.jobFailToProcess(job, err)
	// 	return
	// }

	//TODO: process job
	job.Status = models.StatusCompleted
}

func (w *cvEvaluatorWorker) jobFailToProcess(job *models.JobItem, err error) {
	fmt.Printf("job with id %s failed to prcess: %s\n", job.Id, err.Error())
	job.Status = models.StatusFailed
}

func (w *cvEvaluatorWorker) buildCvEvaluatorPrompt(jobTitle, extractedCv string, jobDescription, cvRubric []models.ChromaSearchResult) string {
	prompt := "Evaluate this CV for role: " + jobTitle + "\n\n"
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
	prompt += "Return as:\nmatch_rate: <0.0-1.0>\n---\nfeedback: <brief feedback with 2-3 sentences>\n"
	return prompt
}
