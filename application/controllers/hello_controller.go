package controllers

import (
	"context"
	"net/http"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/services"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/models"
)

type IHelloController interface {
	GetHello(context.Context, *http.Request) models.WebResponse
}

type helloController struct {
	helloService services.IHelloService
}

func NewHelloController(helloService services.IHelloService) IHelloController {
	return &helloController{
		helloService: helloService,
	}
}

func (h *helloController) GetHello(ctx context.Context, r *http.Request) models.WebResponse {
	return h.helloService.GetHello(ctx)
}
