package command

import (
	"github/imfropz/go-ddd/internal/application/common"

	"github.com/google/uuid"
)

type UpdateProfileCommand struct {
	Id              uuid.UUID
	Name            string
	Email           string
	CurrentPassword string
	NewPassword     string
}

type UpdateProfileCommandResult struct {
	Result *common.UserResult
}
