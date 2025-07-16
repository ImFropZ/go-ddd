package criteria

import "github.com/google/uuid"

type TodoCriteria struct {
	UserId    uuid.UUID
	Completed *bool
}

func (todo *TodoCriteria) WithUserId(id uuid.UUID) *TodoCriteria {
	todo.UserId = id
	return todo
}

func (todo *TodoCriteria) WithCompleted(completed *bool) *TodoCriteria {
	todo.Completed = completed
	return todo
}
