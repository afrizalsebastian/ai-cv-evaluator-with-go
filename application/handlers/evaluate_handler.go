package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/services"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/models"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/worker"
	"github.com/google/uuid"
)

type IEvaluateHandler interface {
	EvaluateFile(w http.ResponseWriter, r *http.Request)
}

type evaluateHandler struct {
	worker   worker.ICvEvaluatorWorker
	jobStore services.IJobStore
}

func NewEvaluateHandler(worker worker.ICvEvaluatorWorker, jobStore services.IJobStore) IEvaluateHandler {
	return &evaluateHandler{
		worker:   worker,
		jobStore: jobStore,
	}
}

func (e *evaluateHandler) EvaluateFile(w http.ResponseWriter, r *http.Request) {
	var requestBody models.EvaluateRequest
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &requestBody); err != nil {
		fmt.Printf("error when parse request body: %s", err.Error())
		resp := models.CreateWebResponse(err.Error(), http.StatusBadRequest, nil)
		respJson, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respJson)
	}

	jobId := uuid.New().String()
	jobItem := &models.JobItem{
		Id:       jobId,
		JobTitle: requestBody.JobTitle,
		FileId:   requestBody.FileId,
		Status:   models.StatusQueued,
	}

	e.jobStore.Set(jobId, jobItem)

	// enqueue
	if err := e.worker.Enqueue(jobItem); err != nil {
		fmt.Printf("error when enqueue job: %s", err.Error())
		resp := models.CreateWebResponse(err.Error(), http.StatusInternalServerError, nil)
		respJson, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respJson)
	}

	data := map[string]interface{}{
		"job_id": jobId,
		"status": jobItem.Status,
	}
	resp := models.CreateWebResponse("Success", http.StatusOK, data)
	b, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)
	w.Write(b)
}
