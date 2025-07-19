//go:generate mockgen -source=publisher.go -destination=../mocks/event_publisher_mock.go -package=mocks

package event

type EventPublisher interface {
	Publish(topic string, event interface{}) error
	PublishWithKey(topic string, key []byte, event interface{}) error
	Close() error
}
