package criteria

import "github.com/google/uuid"

type UserCriteria struct {
	Id   uuid.UUID
	Name *string
}

func (user *UserCriteria) WithId(id uuid.UUID) *UserCriteria {
	user.Id = id
	return user
}

func (user *UserCriteria) WithName(name *string) *UserCriteria {
	user.Name = name
	return user
}
