package kafka

import (
	"context"
	"fmt"
	"github/imfropz/go-ddd/internal/domain/event"
	"log/slog"
	"sync"

	"github.com/IBM/sarama"
)

type SaramaConsumer struct {
	consumer sarama.ConsumerGroup
	config   *SaramaConfig
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewSaramaConsumer(config *SaramaConfig, groupId string) (*SaramaConsumer, error) {
	kafkaConfig, err := createSaramaConfig(config)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumerGroup(config.Brokers, groupId, kafkaConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &SaramaConsumer{
		consumer: consumer,
		config:   config,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

func (c *SaramaConsumer) Consume(topics []string, handler event.EventHandler) error {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		consumerHandler := &consumerGroupHandler{
			handler: handler,
		}

		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				if err := c.consumer.Consume(c.ctx, topics, consumerHandler); err != nil {
					slog.Error(fmt.Sprintf("Error from consumer: %v", err))
				}
			}
		}
	}()
	return nil
}

func (c *SaramaConsumer) Close() error {
	c.cancel()
	c.wg.Wait()
	return c.consumer.Close()
}

type consumerGroupHandler struct {
	handler event.EventHandler
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := h.handler.Handle(message.Topic, message.Key, message.Value); err != nil {
			slog.Error(fmt.Sprintf("Error handling consumer: %v", err))
			continue
		}
		session.MarkMessage(message, "")
	}
	return nil
}

var _ event.EventConsumer = (*SaramaConsumer)(nil)
