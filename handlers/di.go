package handlers

import (
	"fmt"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/controllers"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/services"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/bootstrap"
)

type CvEvaluatorServiceController struct {
	Hello controllers.IHelloController
}

func initDI(app *bootstrap.Application) *CvEvaluatorServiceController {
	init := &CvEvaluatorServiceController{
		Hello: hello(app),
	}

	return init
}

func hello(_ *bootstrap.Application) controllers.IHelloController {
	helloService := services.NewHelloService()
	helloController := controllers.NewHelloController(helloService)
	fmt.Println(helloController)
	return helloController
}
