package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/imfropz/go-ddd/internal/application/command"
	"github/imfropz/go-ddd/internal/application/interfaces"
	"github/imfropz/go-ddd/internal/domain/entity"
	"log/slog"
	"os"
)

type NotificationEventHandler struct {
	notificationService interfaces.NotificationService
}

func NewNotificationEventHandler(notificationService interfaces.NotificationService) *NotificationEventHandler {
	return &NotificationEventHandler{
		notificationService: notificationService,
	}
}

func (handler *NotificationEventHandler) Handle(topic string, key, value []byte) error {
	switch topic {
	case entity.RESET_PASSWORD:
		var event entity.ResetPasswordEvent
		if err := json.Unmarshal(value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal reset password event: %v", err)
		}
		return handler.handleResetPassword(event)
	default:
		return fmt.Errorf("unknown topic: %s", topic)
	}
}

func (handler *NotificationEventHandler) handleResetPassword(event entity.ResetPasswordEvent) error {
	fromEmail := os.Getenv("FROM_EMAIL")
	if fromEmail == "" {
		slog.Error("missing FROM_EMAIL enviorment variable")
		return errors.New("missing FROM_EMAIL enviorment variable")
	}

	// TODO: Update the html body template
	handler.notificationService.SendEmail(&command.SendEmailCommand{
		FromEmail: fromEmail,
		ToEmails:  []string{event.Email},
		Subject:   "Reset Password - Buon18",
		HtmlBody:  `<h1>Hello World</h1> <p>This is a test email.</p>`,
	})
	return nil
}
