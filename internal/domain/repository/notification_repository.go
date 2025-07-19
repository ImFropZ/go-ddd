//go:generate mockgen -source=notification_repository.go -destination=../mocks/notification_repository_mock.go -package=mocks

package repository

type NotificationRepository interface {
	SendToEmail(fromEmail string, toEmail []string, message string) error
}
