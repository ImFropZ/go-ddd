package common

import (
	"time"

	"github.com/google/uuid"
)

type UserResult struct {
	Id        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
