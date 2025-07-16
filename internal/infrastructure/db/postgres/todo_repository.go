package postgres

import (
	"github/imfropz/go-ddd/internal/domain/criteria"
	"github/imfropz/go-ddd/internal/domain/entity"
	"github/imfropz/go-ddd/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormTodoRepository struct {
	db *gorm.DB
}

func NewGormTodoRepository(db *gorm.DB) repository.TodoRepository {
	return &GormTodoRepository{db: db}
}

func (repo *GormTodoRepository) Create(todo *entity.ValidatedTodo) (*entity.Todo, error) {
	dbTodo := toDBTodo(todo)

	if err := repo.db.Create(dbTodo).Error; err != nil {
		return nil, err
	}

	return repo.FindById(dbTodo.Id)
}

func (repo *GormTodoRepository) FindById(id uuid.UUID) (*entity.Todo, error) {
	var dbTodo Todo
	if err := repo.db.First(&dbTodo, id).Error; err != nil {
		return nil, err
	}

	return fromDBTodo(&dbTodo), nil
}

func (repo *GormTodoRepository) FindAll(todoCriteria *criteria.TodoCriteria) ([]*entity.Todo, error) {
	query := repo.db.Model(&Todo{})
	if todoCriteria != nil {
		if todoCriteria.UserId != uuid.Nil {
			query = query.Where("user_id = ?", todoCriteria.UserId)
		}
		if todoCriteria.Completed != nil {
			query = query.Where("completed = ?", todoCriteria.Completed)
		}
	}

	var dbTodos []Todo
	if err := query.Find(&dbTodos).Error; err != nil {
		return nil, err
	}

	var todos = make([]*entity.Todo, len(dbTodos))
	for i, dbTodo := range dbTodos {
		todos[i] = fromDBTodo(&dbTodo)
	}

	return todos, nil
}

func (repo *GormTodoRepository) Update(todo *entity.ValidatedTodo) (*entity.Todo, error) {
	dbTodo := toDBTodo(todo)

	if err := repo.db.Model(&Todo{}).Where("id = ?", dbTodo.Id).Updates(dbTodo).Error; err != nil {
		return nil, err
	}

	return repo.FindById(dbTodo.Id)
}

func (repo *GormTodoRepository) Delete(id uuid.UUID) error {
	return repo.db.Delete(&Todo{}, id).Error
}
