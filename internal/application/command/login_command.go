package command

import "github/imfropz/go-ddd/internal/application/common"

type LoginCommand struct {
	Email    string
	Password string
}

type LoginCommandResult struct {
	Result *common.UserResult
}
