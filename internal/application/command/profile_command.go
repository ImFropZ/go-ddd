package command

import "github/imfropz/go-ddd/internal/application/common"

type ProfileCommand struct {
	Email string
}

type ProfileCommandResult struct {
	Result *common.UserResult
}
