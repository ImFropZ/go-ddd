package request

import (
	"encoding/json"
	"github/imfropz/go-ddd/internal/application/command"
	"io"
	"net/http"
)

type CreateUserRequest struct {
	Name string `json:"name"`
}

func NewCreateUserRequest(r *http.Request) (*CreateUserRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var req CreateUserRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (req *CreateUserRequest) ToCreateUserCommand() *command.CreateUserCommand {
	return &command.CreateUserCommand{
		Name: req.Name,
	}
}
