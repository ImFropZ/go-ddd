package main

import (
	"fmt"
	"github/imfropz/go-ddd/internal/application/service"
	"github/imfropz/go-ddd/internal/infrastructure/db/postgres"
	"github/imfropz/go-ddd/internal/interface/api/rest"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func main() {
	db, err := postgres.NewConnection()
	if err != nil {
		panic("unable connect to database")
	}

	databaseMigration(db)

	userRepository := postgres.NewGormUserRepository(db)
	userService := service.NewUserService(userRepository)

	r := mux.NewRouter()
	rest.NewUserController(r, userService)

	slog.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error(fmt.Sprintf("Failed to start server: %s", err))
	}
}

func databaseMigration(db *gorm.DB) {
	db.AutoMigrate(&postgres.Todo{}, &postgres.User{})
}
