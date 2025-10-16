package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/kafka"
)

type IKafkaProducer interface {
	PublishMessage(ctx context.Context, topic string, key, message interface{}) error
}

type kafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer(producer *kafka.Producer) IKafkaProducer {
	return &kafkaProducer{
		producer: producer,
	}
}

func (k *kafkaProducer) PublishMessage(ctx context.Context, topic string, key, message interface{}) error {
	if k.producer == nil {
		log.Println("kafka not initialized")
		return errors.New("kafka not initialized")
	}

	byteKey, _ := json.Marshal(key)
	byteMessage, _ := json.Marshal(message)

	_, _, err := k.producer.Publish(topic, byteKey, byteMessage)

	if err != nil {
		log.Println("kafka failed to publish message")
		return err
	}

	return nil
}
