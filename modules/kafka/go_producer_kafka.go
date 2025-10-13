package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(brokers []string, options ...ClientOption) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 3
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Timeout = 10 * time.Second
	config.Version = sarama.V3_4_0_0

	for _, option := range options {
		option(config)
	}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer}, nil
}

// Options
func WithSASLAuth(username, password string) ClientOption {
	return func(config *sarama.Config) {
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
	}
}

func WithSASLHandshakeVersion(handshake bool) ClientOption {
	return func(config *sarama.Config) {
		config.Net.SASL.Handshake = handshake
	}
}

func WithSASLMechanism(mechanism sarama.SASLMechanism) ClientOption {
	return func(config *sarama.Config) {
		config.Net.SASL.Mechanism = mechanism
	}
}

func WithSASLEnable(enable bool) ClientOption {
	return func(config *sarama.Config) {
		config.Net.SASL.Enable = enable
	}
}

func WithKafkaVersion(version sarama.KafkaVersion) ClientOption {
	return func(config *sarama.Config) {
		config.Version = version
	}
}

func WithProducerCompression(codec sarama.CompressionCodec) ClientOption {
	return func(config *sarama.Config) {
		config.Producer.Compression = codec
	}
}

func WithProducerFlush(frequency time.Duration) ClientOption {
	return func(config *sarama.Config) {
		config.Producer.Flush.Frequency = frequency
	}
}

func (p *Producer) Publish(topic string, key, value []byte) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	return p.producer.SendMessage(msg)
}

func (p *Producer) PublishWithHeaders(topic string, key, value []byte, headers map[string]string) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic:   topic,
		Key:     sarama.ByteEncoder(key),
		Value:   sarama.ByteEncoder(value),
		Headers: convertHeaders(headers),
	}
	return p.producer.SendMessage(msg)
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

func convertHeaders(headers map[string]string) []sarama.RecordHeader {
	result := make([]sarama.RecordHeader, 0, len(headers))
	for k, v := range headers {
		result = append(result, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}
	return result
}
