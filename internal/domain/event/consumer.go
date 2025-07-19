//go:generate mockgen -source=consumer.go -destination=../mocks/event_consumer_mock.go -package=mocks

package event

type EventConsumer interface {
	Consume(topics []string, handler EventHandler) error
	Close() error
}

type EventHandler interface {
	Handle(topic string, key, value []byte) error
}
