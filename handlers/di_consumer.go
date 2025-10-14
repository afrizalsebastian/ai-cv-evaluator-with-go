package handlers

import (
	controller_consumer "github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/controllers/consumer"
	service_consumer "github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/services/consumer"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/bootstrap"
)

type ConsumerController struct {
	CvEvaluatorConsumer controller_consumer.ICvEvaluatorControllerConsumer
}

func initDIConsumer(app *bootstrap.Application) *ConsumerController {
	initDi := &ConsumerController{
		CvEvaluatorConsumer: cvEvaluatorConsumer(app),
	}

	return initDi
}

func cvEvaluatorConsumer(app *bootstrap.Application) controller_consumer.ICvEvaluatorControllerConsumer {
	cvEvaluatorServiceConsumer := service_consumer.NewCvEvaluatorConsumerService(app.GeminiClient, app.ChromaClient, app.Ingest)
	cvEvaluatorControllerConsumer := controller_consumer.NewCvEvaluatorConsumer(cvEvaluatorServiceConsumer)
	return cvEvaluatorControllerConsumer
}
