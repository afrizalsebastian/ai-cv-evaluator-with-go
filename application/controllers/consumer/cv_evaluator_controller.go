package controller_consumer

import (
	"context"
	"log"

	service_consumer "github.com/afrizalsebastian/ai-cv-evaluator-with-go/application/services/consumer"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/kafka"
)

type cvEvaluatorControllerConsumer struct {
	cvEvaluatorServiceConsumer service_consumer.ICvEvaluatorConsumerService
}

func NewCvEvaluatorConsumer() *cvEvaluatorControllerConsumer {
	return &cvEvaluatorControllerConsumer{}
}

func (c *cvEvaluatorControllerConsumer) ProcessMessage(ctx context.Context, msg *kafka.Message) error {
	request := msg.Value
	jobId := string(request)

	log.Println("running job with id", jobId)
	return c.cvEvaluatorServiceConsumer.RunningJob(ctx, jobId)
}
