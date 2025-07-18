package postgres

import "github/imfropz/go-ddd/internal/domain/entity"

func toDBUser(user *entity.ValidatedUser) *User {
	u := &User{
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	u.Id = user.Id

	return u
}

func fromDBUser(dbUser *User) *entity.User {
	u := &entity.User{
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
	u.Id = dbUser.Id

	return u
}
