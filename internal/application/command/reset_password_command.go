package command

import "github/imfropz/go-ddd/internal/application/common"

type ResetPasswordCommand struct {
	Email string
}

type ResetPasswordCommandResult struct {
	Result *common.UserResult
}

type ResetPasswordWithTokenCommand struct {
	Token       string
	NewPassword string
}

type ResetPasswordWithTokenCommandResult struct {
	Result *common.UserResult
}
