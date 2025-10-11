package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/api"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/controllers"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/bootstrap"
)

type Server struct {
	HelloController          controllers.IHelloController
	UploadDocumentController controllers.IUploadDocumentController
	EvaluateController       controllers.IJobController
}

func NewServer(app *bootstrap.Application) (*Server, error) {
	di := initDI(app)
	server := &Server{
		HelloController:          di.Hello,
		UploadDocumentController: di.UploadDocument,
		EvaluateController:       di.Evaluate,
	}

	return server, nil
}

func (s *Server) GetHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp := s.HelloController.GetHello(ctx, r)
	api.WriteJSONResponse(w, resp.Status, resp)
}

func (s *Server) PostUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp := s.UploadDocumentController.PostDocument(ctx, r)
	api.WriteJSONResponse(w, resp.Status, resp)
}

func (s *Server) PostEvaluate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp := s.EvaluateController.EnqueueJob(ctx, r)
	api.WriteJSONResponse(w, resp.Status, resp)
}

func (s *Server) GetResultJobId(w http.ResponseWriter, r *http.Request, jobId string) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp := s.EvaluateController.ResultJob(ctx, r, jobId)
	api.WriteJSONResponse(w, resp.Status, resp)
}
