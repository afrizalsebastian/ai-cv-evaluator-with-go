package services

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/api"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/domain/models"
	jobstore "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/job-store"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/worker"
	"github.com/google/uuid"
)

type IJobService interface {
	EnqueueJob(context.Context, *models.EvaluateRequest) api.WebResponse
	ResultJob(context.Context, string) api.WebResponse
}

type jobService struct {
	worker   worker.ICvEvaluatorWorker
	jobStore jobstore.IJobStore
}

func NewEvaluateServce(worker worker.ICvEvaluatorWorker, jobStore jobstore.IJobStore) IJobService {
	return &jobService{
		worker:   worker,
		jobStore: jobStore,
	}
}

func (e *jobService) EnqueueJob(ctx context.Context, request *models.EvaluateRequest) api.WebResponse {
	jobId := uuid.New().String()
	jobItem := &models.JobItem{
		Id:       jobId,
		JobTitle: request.JobTitle,
		FileId:   request.FileId,
		Status:   models.StatusQueued,
	}

	e.jobStore.Set(jobId, jobItem)

	// enqueue to worker
	// TODO: maybe can change to kafka publish
	if err := e.worker.Enqueue(jobItem); err != nil {
		log.Println("error to enqueue the job")

		if errors.Is(err, worker.ErrQueueIsFull) {
			return api.CreateWebResponse("queue is full", http.StatusBadRequest, nil, nil)
		}

		return api.CreateWebResponse("internal server error", http.StatusInternalServerError, nil, nil)
	}

	resp := &models.EvaluateResponse{
		JobId:  jobId,
		Status: string(jobItem.Status),
	}

	return api.CreateWebResponse("Success to enqueue the job", http.StatusOK, resp, nil)
}

func (e *jobService) ResultJob(ctx context.Context, jobId string) api.WebResponse {
	jobItem, ok := e.jobStore.Get(jobId)
	if !ok {
		return api.CreateWebResponse("Job Not Found", http.StatusNotFound, nil, nil)
	}

	return api.CreateWebResponse("Success", http.StatusOK, jobItem, nil)
}
