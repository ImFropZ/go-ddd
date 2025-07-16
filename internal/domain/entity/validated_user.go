package entity

type ValidatedUser struct {
	User
	isValidated bool
}

func (vu *ValidatedUser) IsValid() bool {
	return vu.isValidated
}

func NewValidatedUser(user *User) (*ValidatedUser, error) {
	if err := user.validate(); err != nil {
		return nil, err
	}

	return &ValidatedUser{
		User:        *user,
		isValidated: true,
	}, nil
}
