package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

func (u *User) validate() error {
	if u.Name == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

func NewUser(name string) *User {
	return &User{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
}

func (u *User) UpdateName(name string) error {
	u.Name = name
	u.UpdatedAt = time.Now()

	return u.validate()
}
