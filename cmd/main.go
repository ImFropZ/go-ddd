package main

import (
	"fmt"
	"github/imfropz/go-ddd/internal/application/handler"
	"github/imfropz/go-ddd/internal/application/service"
	"github/imfropz/go-ddd/internal/domain/entity"
	"github/imfropz/go-ddd/internal/infrastructure/db/postgres"
	"github/imfropz/go-ddd/internal/infrastructure/db/valkey"
	"github/imfropz/go-ddd/internal/infrastructure/gmail"
	"github/imfropz/go-ddd/internal/infrastructure/kafka"
	"github/imfropz/go-ddd/internal/interface/api"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func main() {
	db, err := postgres.NewConnection()
	if err != nil {
		panic("unable connect to database")
	}

	databaseMigration(db)

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" {
		panic("missing any smtp configuration: SMTP_HOST, SMTP_PORT, SMTP_USERNAME, and SMTP_PASSWORD")
	}

	mail := gmail.NewGmailMail(gmail.SMTPConfig{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
	})

	userRepository := postgres.NewGormUserRepository(db)

	consumer, err := kafka.NewSaramaConsumer(&kafka.SaramaConfig{
		Brokers:  []string{"localhost:9092"},
		Version:  "3.9.1",
		ClientID: "user-service",
	}, "notification-service-group")
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create kafka consumer: %v", err))
		return
	}
	defer consumer.Close()

	valkeyRepository, err := valkey.NewValkeyRepository()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create valkey repository: %v", err))
		return
	}

	notificationService := service.NewNotificationService(mail)

	notificationHandler := handler.NewNotificationEventHandler(notificationService)

	topics := []string{entity.RESET_PASSWORD}
	if err := consumer.Consume(topics, notificationHandler); err != nil {
		slog.Error(fmt.Sprintf("Failed to start consumer: %v", err))
	}

	userProducer, err := kafka.NewSaramaProducer(&kafka.SaramaConfig{
		Brokers:  []string{"localhost:9092"},
		Version:  "3.9.1",
		ClientID: "user-service",
	})
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create kafka producer: %v", err))
		return
	}

	authenticateService := service.NewAuthenticateService(userProducer, valkeyRepository, userRepository)

	r := mux.NewRouter()
	api.NewAuthenticateController(r, authenticateService, userRepository)

	slog.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error(fmt.Sprintf("Failed to start server: %s", err))
	}
}

func databaseMigration(db *gorm.DB) {
	db.AutoMigrate(&postgres.User{})
}
