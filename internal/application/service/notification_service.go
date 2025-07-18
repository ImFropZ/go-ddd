package service

import (
	"github/imfropz/go-ddd/internal/application/command"
	"github/imfropz/go-ddd/internal/infrastructure/gmail"
)

type NotificationService struct {
	mail *gmail.GmailMail
}

func NewNotificationService(mail *gmail.GmailMail) *NotificationService {
	return &NotificationService{
		mail: mail,
	}
}

func (service *NotificationService) SendEmail(sendEmailCommand *command.SendEmailCommand) error {
	return service.mail.SendToEmail(sendEmailCommand.FromEmail, sendEmailCommand.ToEmails, sendEmailCommand.Subject, sendEmailCommand.HtmlBody)
}
