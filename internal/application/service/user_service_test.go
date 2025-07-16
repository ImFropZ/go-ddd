package service_test

import (
	"errors"
	"github/imfropz/go-ddd/internal/application/command"
	"github/imfropz/go-ddd/internal/application/service"
	"github/imfropz/go-ddd/internal/domain/criteria"
	"github/imfropz/go-ddd/internal/domain/entity"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockUserRepository struct {
	users []*entity.ValidatedUser
}

func (mock *MockUserRepository) Create(user *entity.ValidatedUser) (*entity.User, error) {
	mock.users = append(mock.users, user)
	return &user.User, nil
}

func (mock *MockUserRepository) FindById(id uuid.UUID) (*entity.User, error) {
	for _, user := range mock.users {
		if user.Id == id {
			return &user.User, nil
		}
	}
	return nil, errors.New("user not found")
}

func (mock *MockUserRepository) FindAll(userCriteria *criteria.UserCriteria) ([]*entity.User, error) {
	if userCriteria == nil {
		users := make([]*entity.User, 0)
		for _, user := range mock.users {
			users = append(users, &user.User)
		}
		return users, nil
	}

	users := make([]*entity.User, 0)
	for _, user := range mock.users {
		checkCount := 0
		validCount := 0

		if userCriteria.Id != uuid.Nil {
			checkCount++
			if userCriteria.Id == user.Id {
				validCount++
			}
		}
		if userCriteria.Name != nil {
			checkCount++
			if *userCriteria.Name == user.Name {
				validCount++
			}
		}

		if checkCount == validCount {
			users = append(users, &user.User)
		}
	}
	return users, nil
}

func (mock *MockUserRepository) Update(user *entity.ValidatedUser) (*entity.User, error) {
	for i, muser := range mock.users {
		if muser.Id == user.Id {
			mock.users[i] = user
			return &user.User, nil
		}
	}
	return nil, errors.New("user not found for update")
}

func (mock *MockUserRepository) Delete(id uuid.UUID) error {
	for i, user := range mock.users {
		if user.Id == id {
			mock.users = append(mock.users[:i], mock.users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found for delete")
}

func TestUserService_CreateUser(t *testing.T) {
	userRepo := &MockUserRepository{}
	service := service.NewUserService(userRepo)

	user := entity.NewUser("John Doe")
	userCommand := getCreateUserCommand(user)

	_, err := service.CreateUser(userCommand)
	assert.Nil(t, err)
	assert.Len(t, userRepo.users, 1)
}

func TestUserService_FindAllUsers(t *testing.T) {
	userRepo := &MockUserRepository{}
	service := service.NewUserService(userRepo)

	_, _ = service.CreateUser(getCreateUserCommand(entity.NewUser("John Doe")))
	_, _ = service.CreateUser(getCreateUserCommand(entity.NewUser("Jane Doe")))
	_, _ = service.CreateUser(getCreateUserCommand(entity.NewUser("Baby Doe")))

	result, err := service.FindAllUsers(nil)
	assert.Nil(t, err)
	assert.Len(t, result.Result, 3)
}

func TestUserService_FindUserById(t *testing.T) {
	userRepo := &MockUserRepository{}
	service := service.NewUserService(userRepo)

	createUserResult, err := service.CreateUser(getCreateUserCommand(entity.NewUser("John Doe")))
	assert.Nil(t, err)

	userResult, err := service.FindUserById(createUserResult.Result.Id)
	assert.Nil(t, err)
	assert.Equal(t, userResult.Result, createUserResult.Result)
}

func getCreateUserCommand(user *entity.User) *command.CreateUserCommand {
	return &command.CreateUserCommand{
		Name: user.Name,
	}
}
