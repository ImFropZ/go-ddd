package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	RESET_PASSWORD = "reset-password"
)

type User struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     string
	Password  string
}

func (u *User) validate() error {
	if u.Name == "" {
		return errors.New("name must not be empty")
	}
	if u.Email == "" {
		return errors.New("email must not be empty")
	}
	if u.Password == "" {
		return errors.New("password must not be empty")
	}

	return nil
}

func NewUser(name string, email string, password string) *User {
	return &User{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Email:     email,
		Password:  password,
	}
}

func (u *User) UpdateName(name string) error {
	u.Name = name
	u.UpdatedAt = time.Now()

	return u.validate()
}

func (u *User) UpdateEmail(email string) error {
	u.Email = email
	u.UpdatedAt = time.Now()

	return u.validate()
}

func (u *User) UpdatePassword(password string) error {
	u.Password = password
	u.UpdatedAt = time.Now()

	return u.validate()
}

type ResetPasswordEvent struct {
	Email string
	Token string
	Exp   time.Time
}
