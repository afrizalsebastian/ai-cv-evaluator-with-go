package router

import (
	"net/http"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/handlers"
	"github.com/gorilla/mux"
)

type IRouter interface {
	Router() http.Handler
}

type router struct {
	uploadHandler   handlers.IUploadFileHandler
	evaluateHandler handlers.IEvaluateHandler
	resultHandler   handlers.IResultHandler
	chromaHandler   handlers.IChromaHandler
}

func NewRouter(uploadHandler handlers.IUploadFileHandler,
	evaluateHandler handlers.IEvaluateHandler,
	resultHandler handlers.IResultHandler,
	chromaHandler handlers.IChromaHandler) IRouter {
	return &router{
		uploadHandler:   uploadHandler,
		evaluateHandler: evaluateHandler,
		resultHandler:   resultHandler,
		chromaHandler:   chromaHandler,
	}
}

func (app *router) Router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/upload", app.uploadHandler.Upload).Methods("POST")
	r.HandleFunc("/evaluate", app.evaluateHandler.EvaluateFile).Methods("POST")
	r.HandleFunc("/result/{jobId}", app.resultHandler.GetResult).Methods("GET")
	r.HandleFunc("/chroma/upsert", app.chromaHandler.Upsert).Methods("PUT")
	return r
}
