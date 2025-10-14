package handlers

import (
	controller_consumer "github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/controllers/consumer"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/bootstrap"
)

type Consumer struct {
	CvEvaluatorConsumer controller_consumer.ICvEvaluatorControllerConsumer
}

func NewConsumer(app *bootstrap.Application) (*Consumer, error) {
	di := initDIConsumer(app)
	consumer := &Consumer{
		CvEvaluatorConsumer: di.CvEvaluatorConsumer,
	}

	return consumer, nil
}
