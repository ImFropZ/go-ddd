package postgres

import "github/imfropz/go-ddd/internal/domain/entity"

func toDBTodo(todo *entity.ValidatedTodo) *Todo {
	t := &Todo{
		Title:     todo.Title,
		Duration:  todo.Duration,
		Completed: todo.Completed,
		UserId:    todo.User.Id,
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}
	t.Id = todo.Id

	return t
}

func fromDBTodo(dbTodo *Todo) *entity.Todo {
	u := &entity.User{
		Id:        dbTodo.User.Id,
		Name:      dbTodo.User.Name,
		CreatedAt: dbTodo.User.CreatedAt,
		UpdatedAt: dbTodo.User.UpdatedAt,
	}

	t := &entity.Todo{
		Title:     dbTodo.Title,
		Duration:  dbTodo.Duration,
		Completed: dbTodo.Completed,
		User:      *u,
		CreatedAt: dbTodo.CreatedAt,
		UpdatedAt: dbTodo.UpdatedAt,
	}
	t.Id = dbTodo.Id

	return t
}
