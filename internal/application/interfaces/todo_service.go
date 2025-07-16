package interfaces

import (
	"github/imfropz/go-ddd/internal/application/command"
	"github/imfropz/go-ddd/internal/application/query"
	"github/imfropz/go-ddd/internal/domain/criteria"

	"github.com/google/uuid"
)

type TodoService interface {
	CreateTodo(todoCommand *command.CreateTodoCommand) (*command.CreateTodoCommandResult, error)
	FindAllTodos(todoCriteria *criteria.TodoCriteria) (*query.TodoQueryListResult, error)
	FindAllTodosByUserId(userId uuid.UUID, todoCriteria *criteria.TodoCriteria) (*query.TodoQueryListResult, error)
	DeleteTodo(id uuid.UUID) error
}
