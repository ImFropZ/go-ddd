package service

import (
	"github/imfropz/go-ddd/internal/application/command"
	"github/imfropz/go-ddd/internal/application/mapper"
	"github/imfropz/go-ddd/internal/application/query"
	"github/imfropz/go-ddd/internal/domain/criteria"
	"github/imfropz/go-ddd/internal/domain/entity"
	"github/imfropz/go-ddd/internal/domain/repository"

	"github.com/google/uuid"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (service *UserService) CreateUser(userCommand *command.CreateUserCommand) (*command.CreateUserCommandResult, error) {
	userEntity := entity.NewUser(userCommand.Name)

	validatedUser, err := entity.NewValidatedUser(userEntity)
	if err != nil {
		return nil, err
	}

	_, err = service.userRepository.Create(validatedUser)
	if err != nil {
		return nil, err
	}

	result := command.CreateUserCommandResult{
		Result: mapper.NewUserResultFromValidatedEntity(validatedUser),
	}

	return &result, nil
}

func (service *UserService) FindAllUsers(userCriteria *criteria.UserCriteria) (*query.UserQueryListResult, error) {
	users, err := service.userRepository.FindAll(userCriteria)
	if err != nil {
		return nil, err
	}

	var queryListResult query.UserQueryListResult
	for _, user := range users {
		queryListResult.Result = append(queryListResult.Result, mapper.NewUserResultFromEntity(user))
	}

	return &queryListResult, nil
}

func (service *UserService) FindUserById(id uuid.UUID) (*query.UserQueryResult, error) {
	user, err := service.userRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	queryResult := query.UserQueryResult{
		Result: mapper.NewUserResultFromEntity(user),
	}

	return &queryResult, nil
}
