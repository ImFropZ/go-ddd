package repository

import (
	"github/imfropz/go-ddd/internal/domain/criteria"
	"github/imfropz/go-ddd/internal/domain/entity"

	"github.com/google/uuid"
)

type TodoRepository interface {
	Create(todo *entity.ValidatedTodo) (*entity.Todo, error)
	FindById(id uuid.UUID) (*entity.Todo, error)
	FindAll(todoCriteria *criteria.TodoCriteria) ([]*entity.Todo, error)
	Update(todo *entity.ValidatedTodo) (*entity.Todo, error)
	Delete(id uuid.UUID) error
}
