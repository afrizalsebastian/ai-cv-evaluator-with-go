package handlers

import (
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/controllers"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/services"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/bootstrap"
)

type CvEvaluatorServiceController struct {
	Hello          controllers.IHelloController
	UploadDocument controllers.IUploadDocumentController
	Evaluate       controllers.IJobController
}

func initDI(app *bootstrap.Application) *CvEvaluatorServiceController {
	init := &CvEvaluatorServiceController{
		Hello:          hello(app),
		UploadDocument: uploadDocument(app),
		Evaluate:       evaluate(app),
	}

	return init
}

func hello(_ *bootstrap.Application) controllers.IHelloController {
	helloService := services.NewHelloService()
	helloController := controllers.NewHelloController(helloService)
	return helloController
}

func uploadDocument(_ *bootstrap.Application) controllers.IUploadDocumentController {
	uploadDocumentService := services.NewUploadDocumentService("./uploaded-file")
	uploadDocumentController := controllers.NewUploadDocumenController(uploadDocumentService)
	return uploadDocumentController
}

func evaluate(app *bootstrap.Application) controllers.IJobController {
	evaluateService := services.NewEvaluateServce(app.Worker, app.JobStore)
	evaluateController := controllers.NewEvaluateController(evaluateService)
	return evaluateController
}
