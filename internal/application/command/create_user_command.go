package command

import (
	"github/imfropz/go-ddd/internal/application/common"
)

type CreateUserCommand struct {
	Name string
}

type CreateUserCommandResult struct {
	Result *common.UserResult
}
