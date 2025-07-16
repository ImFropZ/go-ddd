package command

import (
	"github/imfropz/go-ddd/internal/application/common"

	"github.com/google/uuid"
)

type CreateTodoCommand struct {
	Title     string
	Duration  int
	UserId    uuid.UUID
	Completed bool
}

type CreateTodoCommandResult struct {
	Result *common.TodoResult
}
