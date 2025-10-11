package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/api"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/helper"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/services"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/domain/models"
)

type IEvaluateController interface {
	EnqueueJob(ctx context.Context, r *http.Request) api.WebResponse
}

type evaluateController struct {
	evaluateService services.IEvaluateService
}

func NewEvaluateController(
	evaluateService services.IEvaluateService,
) IEvaluateController {
	return &evaluateController{
		evaluateService: evaluateService,
	}
}

func (e *evaluateController) EnqueueJob(ctx context.Context, r *http.Request) api.WebResponse {
	request, err := helper.ParseJSONBodyRequest[models.EvaluateRequest](r)
	if err != nil {
		log.Println("error when parse body request")
		return api.CreateWebResponse("invalid request", http.StatusBadRequest, nil, nil)
	}

	// validation
	if err := helper.ValidateParams(ctx, request); err != nil {
		log.Println("validation error")
		return api.CreateWebResponse("validation error", http.StatusBadRequest, nil, err)
	}

	resp := e.evaluateService.EnqueueJob(ctx, request)
	return resp
}
