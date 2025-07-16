package mapper

import (
	"github/imfropz/go-ddd/internal/application/common"
	"github/imfropz/go-ddd/internal/domain/entity"
)

func NewUserResultFromEntity(user *entity.User) *common.UserResult {
	if user == nil {
		return nil
	}

	return &common.UserResult{
		Id:        user.Id,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func NewUserResultFromValidatedEntity(user *entity.ValidatedUser) *common.UserResult {
	return NewUserResultFromEntity(&user.User)
}
