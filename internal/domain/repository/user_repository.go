package repository

import (
	"github/imfropz/go-ddd/internal/domain/criteria"
	"github/imfropz/go-ddd/internal/domain/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *entity.ValidatedUser) (*entity.User, error)
	FindById(id uuid.UUID) (*entity.User, error)
	FindAll(userCriteria *criteria.UserCriteria) ([]*entity.User, error)
	Update(user *entity.ValidatedUser) (*entity.User, error)
	Delete(id uuid.UUID) error
}
