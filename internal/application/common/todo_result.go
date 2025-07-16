package common

import (
	"time"

	"github.com/google/uuid"
)

type TodoResult struct {
	Id        uuid.UUID
	Title     string
	Duration  int
	Completed bool
	User      *UserResult
	CreatedAt time.Time
	UpdatedAt time.Time
}
