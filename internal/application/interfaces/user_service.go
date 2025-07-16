package interfaces

import (
	"github/imfropz/go-ddd/internal/application/command"
	"github/imfropz/go-ddd/internal/application/query"
	"github/imfropz/go-ddd/internal/domain/criteria"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(userCommand *command.CreateUserCommand) (*command.CreateUserCommandResult, error)
	FindAllUsers(userCriteria *criteria.UserCriteria) (*query.UserQueryListResult, error)
	FindUserById(id uuid.UUID) (*query.UserQueryResult, error)
}
