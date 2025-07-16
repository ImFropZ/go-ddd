package mapper

import (
	"github/imfropz/go-ddd/internal/application/common"
	"github/imfropz/go-ddd/internal/domain/entity"
)

func NewTodoResultFromEntity(todo *entity.Todo) *common.TodoResult {
	if todo == nil {
		return nil
	}

	return &common.TodoResult{
		Id:        todo.Id,
		Title:     todo.Title,
		Duration:  todo.Duration,
		Completed: todo.Completed,
		User:      NewUserResultFromEntity(&todo.User),
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}
}

func NewTodoResultFromValidatedEntity(todo *entity.ValidatedTodo) *common.TodoResult {
	return NewTodoResultFromEntity(&todo.Todo)
}
