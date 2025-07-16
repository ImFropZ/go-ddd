package postgres

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	Id        uuid.UUID `gorm:"primaryKey"`
	Title     string
	Duration  int
	Completed bool
	UserId    uuid.UUID `gorm:"index"`
	User      User      `gorm:"foreignKey:UserId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	Id        uuid.UUID `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
