package entity

type ValidatedTodo struct {
	Todo
	isValidated bool
}

func (vt *ValidatedTodo) IsValid() bool {
	return vt.isValidated
}

func NewValidatedTodo(todo *Todo) (*ValidatedTodo, error) {
	if err := todo.validate(); err != nil {
		return nil, err
	}

	return &ValidatedTodo{
		Todo:        *todo,
		isValidated: true,
	}, nil
}
