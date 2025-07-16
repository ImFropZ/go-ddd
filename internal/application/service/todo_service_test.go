package service_test

import (
	"errors"
	"github/imfropz/go-ddd/internal/application/command"
	"github/imfropz/go-ddd/internal/application/service"
	"github/imfropz/go-ddd/internal/domain/criteria"
	"github/imfropz/go-ddd/internal/domain/entity"
	"github/imfropz/go-ddd/internal/domain/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockTodoRepository struct {
	todos []*entity.ValidatedTodo
}

func (mock *MockTodoRepository) Create(todo *entity.ValidatedTodo) (*entity.Todo, error) {
	mock.todos = append(mock.todos, todo)
	return &todo.Todo, nil
}

func (mock *MockTodoRepository) FindById(id uuid.UUID) (*entity.Todo, error) {
	for _, todo := range mock.todos {
		if todo.Id == id {
			return &todo.Todo, nil
		}
	}
	return nil, errors.New("todo not found")
}

func (mock *MockTodoRepository) FindAll(todoCriteria *criteria.TodoCriteria) ([]*entity.Todo, error) {
	if todoCriteria == nil {
		var todos []*entity.Todo
		for _, todo := range mock.todos {
			todos = append(todos, &todo.Todo)
		}
		return todos, nil
	}

	var todos []*entity.Todo
	for _, todo := range mock.todos {
		checkCount := 0
		validCount := 0

		if todoCriteria.UserId != uuid.Nil {
			checkCount++
			if todoCriteria.UserId == todo.User.Id {
				validCount++
			}
		}

		if todoCriteria.Completed != nil {
			checkCount++
			if *todoCriteria.Completed == todo.Completed {
				validCount++
			}
		}

		if checkCount == validCount {
			todos = append(todos, &todo.Todo)
		}
	}

	return todos, nil
}

func (mock *MockTodoRepository) Update(todo *entity.ValidatedTodo) (*entity.Todo, error) {
	for i, mtodo := range mock.todos {
		if mtodo.Id == todo.Id {
			mock.todos[i] = todo
			return &todo.Todo, nil
		}
	}
	return nil, errors.New("todo not found for update")
}

func (mock *MockTodoRepository) Delete(id uuid.UUID) error {
	for i, todo := range mock.todos {
		if todo.Id == id {
			mock.todos = append(mock.todos[:i], mock.todos[i+1:]...)
			return nil
		}
	}
	return errors.New("todo not found for delete")
}

func TestTodoService_CreateTodo(t *testing.T) {
	todoRepo := &MockTodoRepository{}
	userRepo := &MockUserRepository{}
	service := service.NewTodoService(todoRepo, userRepo)

	user := createPersistedUser(t, userRepo)

	todo := entity.NewTodo("Todo 1", 1, *user, false)
	todoCommand := getCreateTodoCommand(todo)

	_, err := service.CreateTodo(todoCommand)
	assert.Nil(t, err)
	assert.Len(t, todoRepo.todos, 1)
}

func TestTodoService_FindAllTodos(t *testing.T) {
	todoRepo := &MockTodoRepository{}
	userRepo := &MockUserRepository{}
	service := service.NewTodoService(todoRepo, userRepo)

	user := createPersistedUser(t, userRepo)

	_, _ = service.CreateTodo(getCreateTodoCommand(entity.NewTodo("Todo 1", 1, *user, false)))
	_, _ = service.CreateTodo(getCreateTodoCommand(entity.NewTodo("Todo 2", 2, *user, true)))
	_, _ = service.CreateTodo(getCreateTodoCommand(entity.NewTodo("Todo 3", 3, *user, false)))

	result, err := service.FindAllTodos(nil)
	assert.Nil(t, err)
	assert.Len(t, result.Result, 3)
}

func TestTodoService_FindAllTodosByUserId(t *testing.T) {
	todoRepo := &MockTodoRepository{}
	userRepo := &MockUserRepository{}
	service := service.NewTodoService(todoRepo, userRepo)

	user := createPersistedUser(t, userRepo)

	_, _ = service.CreateTodo(getCreateTodoCommand(entity.NewTodo("Todo 1", 1, *user, false)))
	_, _ = service.CreateTodo(getCreateTodoCommand(entity.NewTodo("Todo 2", 2, *user, true)))
	_, _ = service.CreateTodo(getCreateTodoCommand(entity.NewTodo("Todo 3", 3, *user, false)))

	result, err := service.FindAllTodosByUserId(user.Id, nil)
	assert.Nil(t, err)
	assert.Len(t, result.Result, 3)
}

func TestTodoService_FindAllTodosByIncompleted(t *testing.T) {
	todoRepo := &MockTodoRepository{}
	userRepo := &MockUserRepository{}
	service := service.NewTodoService(todoRepo, userRepo)

	user := createPersistedUser(t, userRepo)

	_, _ = service.CreateTodo(getCreateTodoCommand(entity.NewTodo("Todo 1", 1, *user, false)))
	_, _ = service.CreateTodo(getCreateTodoCommand(entity.NewTodo("Todo 2", 2, *user, true)))
	_, _ = service.CreateTodo(getCreateTodoCommand(entity.NewTodo("Todo 3", 3, *user, false)))

	completed := false
	result, err := service.FindAllTodosByUserId(user.Id, &criteria.TodoCriteria{
		Completed: &completed,
	})
	assert.Nil(t, err)
	assert.Len(t, result.Result, 2)
}

func TestTodoService_DeleteTodo(t *testing.T) {
	todoRepo := &MockTodoRepository{}
	userRepo := &MockUserRepository{}
	service := service.NewTodoService(todoRepo, userRepo)

	user := createPersistedUser(t, userRepo)

	_, _ = service.CreateTodo(getCreateTodoCommand(entity.NewTodo("Todo 1", 1, *user, false)))
	todo, err := service.CreateTodo(getCreateTodoCommand(entity.NewTodo("Todo 2", 2, *user, true)))
	_, _ = service.CreateTodo(getCreateTodoCommand(entity.NewTodo("Todo 3", 3, *user, false)))
	assert.Nil(t, err)

	err = service.DeleteTodo(todo.Result.Id)
	assert.Nil(t, err)

	assert.Len(t, todoRepo.todos, 2)
}

func getCreateTodoCommand(todo *entity.Todo) *command.CreateTodoCommand {
	return &command.CreateTodoCommand{
		Title:     todo.Title,
		Duration:  todo.Duration,
		UserId:    todo.User.Id,
		Completed: todo.Completed,
	}
}

func createPersistedUser(t *testing.T, userRepo repository.UserRepository) *entity.ValidatedUser {
	user := entity.NewUser("John Doe")
	validatedUser, err := entity.NewValidatedUser(user)
	require.Nil(t, err)
	_, err = userRepo.Create(validatedUser)
	require.Nil(t, err)
	return validatedUser
}
