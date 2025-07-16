package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
	Duration  int
	Completed bool
	User      User
}

func (t *Todo) validate() error {
	if t.Title == "" {
		return errors.New("title must not be empty")
	}
	if t.Duration <= 0 {
		return errors.New("duration must not be less than or equal to 0")
	}
	return nil
}

func NewTodo(title string, duration int, user ValidatedUser, completed bool) *Todo {
	return &Todo{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     title,
		Duration:  duration,
		User:      user.User,
		Completed: completed,
	}
}

func (t *Todo) UpdateTitle(title string) error {
	t.Title = title
	t.UpdatedAt = time.Now()

	return t.validate()
}

func (t *Todo) UpdateDuration(duration int) error {
	t.Duration = duration
	t.UpdatedAt = time.Now()

	return t.validate()
}
