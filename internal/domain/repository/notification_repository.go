package repository

type NotificationRepository interface {
	SendToEmail(fromEmail string, toEmail []string, message string) error
}
