package service

import (
	"github/imfropz/go-ddd/internal/application/command"
	"github/imfropz/go-ddd/internal/application/mapper"
	"github/imfropz/go-ddd/internal/application/query"
	"github/imfropz/go-ddd/internal/domain/criteria"
	"github/imfropz/go-ddd/internal/domain/entity"
	"github/imfropz/go-ddd/internal/domain/repository"

	"github.com/google/uuid"
)

type TodoService struct {
	todoRepository repository.TodoRepository
	userRepository repository.UserRepository
}

func NewTodoService(todoRepository repository.TodoRepository, userRepository repository.UserRepository) *TodoService {
	return &TodoService{
		todoRepository: todoRepository,
		userRepository: userRepository,
	}
}

func (service *TodoService) CreateTodo(todoCommand *command.CreateTodoCommand) (*command.CreateTodoCommandResult, error) {
	user, err := service.userRepository.FindById(todoCommand.UserId)
	if err != nil {
		return nil, err
	}
	validatedUser, err := entity.NewValidatedUser(user)
	if err != nil {
		return nil, err
	}

	todo := entity.NewTodo(todoCommand.Title, todoCommand.Duration, *validatedUser, todoCommand.Completed)
	validatedTodo, err := entity.NewValidatedTodo(todo)
	if err != nil {
		return nil, err
	}

	_, err = service.todoRepository.Create(validatedTodo)
	if err != nil {
		return nil, err
	}

	result := command.CreateTodoCommandResult{
		Result: mapper.NewTodoResultFromValidatedEntity(validatedTodo),
	}

	return &result, nil
}

func (service *TodoService) FindAllTodos(todoCriteria *criteria.TodoCriteria) (*query.TodoQueryListResult, error) {
	todos, err := service.todoRepository.FindAll(todoCriteria)
	if err != nil {
		return nil, err
	}

	var queryListResult query.TodoQueryListResult
	for _, todo := range todos {
		queryListResult.Result = append(queryListResult.Result, mapper.NewTodoResultFromEntity(todo))
	}

	return &queryListResult, nil
}

func (service *TodoService) FindAllTodosByUserId(userId uuid.UUID, todoCriteria *criteria.TodoCriteria) (*query.TodoQueryListResult, error) {
	if todoCriteria == nil {
		todoCriteria = &criteria.TodoCriteria{}
	}
	return service.FindAllTodos(todoCriteria.WithUserId(userId))
}

func (service *TodoService) DeleteTodo(id uuid.UUID) error {
	return service.todoRepository.Delete(id)
}
