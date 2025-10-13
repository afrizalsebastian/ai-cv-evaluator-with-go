package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

var (
	ErrNoBrokers            = errors.New("no brokers configured")
	ErrNoGroupID            = errors.New("no group ID configured")
	ErrNoTopics             = errors.New("no topics configured")
	ErrNoConsumerController = errors.New("no controller configured")
)

type AckMode int

const (
	AckModeAuto AckMode = iota
	AckModeManual
)

type Message struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       []byte
	Value     []byte
	Timestamp time.Time
	Headers   []*sarama.RecordHeader
}

type ConsumerController interface {
	ProcessMessage(ctx context.Context, msg *Message) error
}

type ConsumerOption func(*consumerConfig)

type consumerConfig struct {
	brokers      []string
	groupID      string
	topics       []string
	controller   ConsumerController
	saramaConfig *sarama.Config
	retryPolicy  RetryPolicy
	ackMode      AckMode
}

type RetryPolicy struct {
	MaxAttempts int
	Backoff     time.Duration
}

// Consumer
type Consumer struct {
	group     sarama.ConsumerGroup
	topics    []string
	handler   *consumerHandler
	closeChan chan struct{}
}

func NewConsumer(opts ...ConsumerOption) (*Consumer, error) {
	config := &consumerConfig{
		saramaConfig: sarama.NewConfig(),
		retryPolicy: RetryPolicy{
			MaxAttempts: 3,
			Backoff:     1 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(config)
	}

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid consumer config: %w", err)
	}

	configureDefaults(config.saramaConfig)

	group, err := sarama.NewConsumerGroup(config.brokers, config.groupID, config.saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &Consumer{
		group:     group,
		topics:    config.topics,
		handler:   newConsumerHandler(config.controller, config.retryPolicy, config.ackMode),
		closeChan: make(chan struct{}),
	}, nil
}

// Run starts the consumer in a blocking manner
func (c *Consumer) Run(ctx context.Context) error {
	defer close(c.closeChan)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-c.group.Errors():
			log.Printf("Consumer group error: %v", err)
		default:
			if err := c.group.Consume(ctx, c.topics, c.handler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return nil
				}
				return fmt.Errorf("consumption error: %w", err)
			}
		}
	}
}

func (c *Consumer) Close() error {
	close(c.closeChan)
	return c.group.Close()
}

// consumerHandler
type consumerHandler struct {
	controller  ConsumerController
	retryPolicy RetryPolicy
	ackMode     AckMode
}

func newConsumerHandler(controller ConsumerController, retryPolicy RetryPolicy, ackMode AckMode) *consumerHandler {
	return &consumerHandler{
		controller:  controller,
		retryPolicy: retryPolicy,
		ackMode:     ackMode,
	}
}

func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer group setup completed")
	return nil
}

func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer group cleanup completed")
	return nil
}

func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		kafkaMsg := &Message{
			Topic:     msg.Topic,
			Partition: msg.Partition,
			Offset:    msg.Offset,
			Key:       msg.Key,
			Value:     msg.Value,
			Timestamp: msg.Timestamp,
			Headers:   msg.Headers,
		}

		if err := h.processWithRetry(session.Context(), kafkaMsg); err != nil {
			log.Printf("Failed to process message after retries: %v", err)
			if h.ackMode == AckModeAuto {
				session.MarkMessage(msg, "")
			}
			continue
		}

		if h.ackMode == AckModeAuto {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}

// Options
func WithBrokers(brokers ...string) ConsumerOption {
	return func(c *consumerConfig) {
		c.brokers = brokers
	}
}

func WithGroupID(groupID string) ConsumerOption {
	return func(c *consumerConfig) {
		c.groupID = groupID
	}
}

func WithTopics(topics ...string) ConsumerOption {
	return func(c *consumerConfig) {
		c.topics = topics
	}
}

func WithConsumerController(controller ConsumerController) ConsumerOption {
	return func(c *consumerConfig) {
		c.controller = controller
	}
}

func WithSaramaConfig(cfg *sarama.Config) ConsumerOption {
	return func(c *consumerConfig) {
		c.saramaConfig = cfg
	}
}

func WithRetryPolicy(maxAttempts int, backoff time.Duration) ConsumerOption {
	return func(c *consumerConfig) {
		c.retryPolicy = RetryPolicy{
			MaxAttempts: maxAttempts,
			Backoff:     backoff,
		}
	}
}

func WithAckMode(mode AckMode) ConsumerOption {
	return func(c *consumerConfig) {
		c.ackMode = mode
	}
}

// internal function
func (h *consumerHandler) processWithRetry(ctx context.Context, msg *Message) error {
	var lastErr error
	for attempt := 1; attempt <= h.retryPolicy.MaxAttempts; attempt++ {
		if err := h.controller.ProcessMessage(ctx, msg); err != nil {
			lastErr = err
			log.Printf("Processing attempt %d failed: %v", attempt, err)
			if attempt < h.retryPolicy.MaxAttempts {
				time.Sleep(h.retryPolicy.Backoff)
			}
			continue
		}
		return nil
	}
	return fmt.Errorf("max retries exceeded (%d): %w", h.retryPolicy.MaxAttempts, lastErr)
}

func validateConfig(config *consumerConfig) error {
	if len(config.brokers) == 0 {
		return ErrNoBrokers
	}
	if config.groupID == "" {
		return ErrNoGroupID
	}
	if len(config.topics) == 0 {
		return ErrNoTopics
	}
	if config.controller == nil {
		return ErrNoConsumerController
	}
	return nil
}

func configureDefaults(config *sarama.Config) {
	if config.Version == (sarama.KafkaVersion{}) {
		config.Version = sarama.V2_5_0_0
	}
	if config.Consumer.Return.Errors != true {
		config.Consumer.Return.Errors = true
	}
	if config.Consumer.Offsets.Initial == 0 {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
}
