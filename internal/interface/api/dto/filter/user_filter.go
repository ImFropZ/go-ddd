package filter

import (
	"github/imfropz/go-ddd/internal/domain/criteria"
	"net/http"

	"github.com/google/uuid"
)

func RequestToUserCriteria(r http.Request) (*criteria.UserCriteria, error) {
	query := r.URL.Query()

	userCriteria := criteria.UserCriteria{}
	if query.Has("id") {
		id, err := uuid.Parse(query.Get("id"))
		if err != nil {
			return nil, err
		}
		userCriteria = *userCriteria.WithId(id)
	}
	if query.Has("name") {
		name := query.Get("name")
		userCriteria = *userCriteria.WithName(&name)
	}

	return &userCriteria, nil
}
