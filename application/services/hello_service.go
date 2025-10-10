package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/models"
)

type IHelloService interface {
	GetHello(context.Context) models.WebResponse
}

type helloService struct{}

func NewHelloService() IHelloService {
	return &helloService{}
}

func (h *helloService) GetHello(ctx context.Context) models.WebResponse {
	fmt.Println("helloService.GetHello")
	return *models.CreateWebResponse("success", http.StatusOK, nil)
}
