package request

import (
	"encoding/json"
	"github/imfropz/go-ddd/internal/application/command"
	"io"
	"net/http"
)

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func NewRegisterRequest(r *http.Request) (*RegisterRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var req RegisterRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (req *RegisterRequest) ToRegisterCommand() *command.RegisterCommand {
	return &command.RegisterCommand{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}
