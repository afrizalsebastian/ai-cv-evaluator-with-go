package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/services"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/models"
	"github.com/gorilla/mux"
)

type IResultHandler interface {
	GetResult(w http.ResponseWriter, r *http.Request)
}

type resultHandler struct {
	jobStore services.IJobStore
}

func NewResultHandler(jobStore services.IJobStore) IResultHandler {
	return &resultHandler{
		jobStore: jobStore,
	}
}

func (e *resultHandler) GetResult(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobId := vars["jobId"]
	job, ok := e.jobStore.Get(jobId)
	if !ok {
		err := errors.New(fmt.Sprintf("job with id %s not found", jobId))
		fmt.Printf("error when get result: %s", err.Error())
		resp := models.CreateWebResponse(err.Error(), http.StatusBadRequest, nil)
		respJson, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respJson)
		return
	}

	jsonResult, _ := json.Marshal(job)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonResult)
}
