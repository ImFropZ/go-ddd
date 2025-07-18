package interfaces

import "github/imfropz/go-ddd/internal/application/command"

type NotificationService interface {
	SendEmail(sendEmailCommand *command.SendEmailCommand) error
}
