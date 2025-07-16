package postgres

import (
	"github/imfropz/go-ddd/internal/domain/criteria"
	"github/imfropz/go-ddd/internal/domain/entity"
	"github/imfropz/go-ddd/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) repository.UserRepository {
	return &GormUserRepository{db: db}
}

func (repo *GormUserRepository) Create(user *entity.ValidatedUser) (*entity.User, error) {
	dbUser := toDBUser(user)

	if err := repo.db.Create(dbUser).Error; err != nil {
		return nil, err
	}

	return repo.FindById(dbUser.Id)
}

func (repo *GormUserRepository) FindById(id uuid.UUID) (*entity.User, error) {
	var dbUser User
	if err := repo.db.First(&dbUser, id).Error; err != nil {
		return nil, err
	}

	return fromDBUser(&dbUser), nil
}

func (repo *GormUserRepository) FindAll(userCriteria *criteria.UserCriteria) ([]*entity.User, error) {
	query := repo.db.Model(&User{})

	if userCriteria != nil {
		if userCriteria.Id != uuid.Nil {
			query = query.Where("id = ?", userCriteria.Id)
		}
		if userCriteria.Name != nil {
			query = query.Where("name = ?", userCriteria.Name)
		}
	}

	var dbUsers []User
	if err := query.Find(&dbUsers).Error; err != nil {
		return nil, err
	}

	users := make([]*entity.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = fromDBUser(&dbUser)
	}

	return users, nil
}

func (repo *GormUserRepository) Update(user *entity.ValidatedUser) (*entity.User, error) {
	dbUser := toDBUser(user)

	if err := repo.db.Model(&User{}).Where("id = ?", dbUser.Id).Updates(dbUser).Error; err != nil {
		return nil, err
	}

	return repo.FindById(dbUser.Id)
}

func (repo *GormUserRepository) Delete(id uuid.UUID) error {
	return repo.db.Delete(&User{}, id).Error
}
