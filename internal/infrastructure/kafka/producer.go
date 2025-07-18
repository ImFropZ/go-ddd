package kafka

import (
	"encoding/json"
	"errors"
	"github/imfropz/go-ddd/internal/domain/event"
	"sync"

	"github.com/IBM/sarama"
)

type SaramaProducer struct {
	producer sarama.SyncProducer
	config   *SaramaConfig
	mu       sync.Mutex
}

type SaramaConfig struct {
	Brokers  []string
	Version  string
	ClientID string
}

func NewSaramaProducer(config *SaramaConfig) (*SaramaProducer, error) {
	kafkaConfig, err := createSaramaConfig(config)
	if err != nil {
		return nil, err
	}

	producer, err := sarama.NewSyncProducer(config.Brokers, kafkaConfig)
	if err != nil {
		return nil, err
	}

	return &SaramaProducer{
		producer: producer,
		config:   config,
	}, nil
}

func createSaramaConfig(config *SaramaConfig) (*sarama.Config, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Retry.Max = 3
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.ClientID = config.ClientID

	version, err := sarama.ParseKafkaVersion(config.Version)
	if err != nil {
		return nil, err
	}
	saramaConfig.Version = version

	return saramaConfig, nil
}

func (p *SaramaProducer) Publish(topic string, event interface{}) error {
	return p.PublishWithKey(topic, nil, event)
}

func (p *SaramaProducer) PublishWithKey(topic string, key []byte, event interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.producer == nil {
		return errors.New("producer is closed")
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(eventBytes),
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *SaramaProducer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.producer != nil {
		err := p.producer.Close()
		p.producer = nil
		return err
	}
	return nil
}

var _ event.EventPublisher = (*SaramaProducer)(nil)
