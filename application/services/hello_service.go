package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/api"
)

type IHelloService interface {
	GetHello(context.Context) api.WebResponse
}

type helloService struct{}

func NewHelloService() IHelloService {
	return &helloService{}
}

func (h *helloService) GetHello(ctx context.Context) api.WebResponse {
	fmt.Println("helloService.GetHello")
	return api.CreateWebResponse("success", http.StatusOK, nil, nil)
}
