package request

import (
	"encoding/json"
	"github/imfropz/go-ddd/internal/application/command"
	"io"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func NewLoginRequest(r *http.Request) (*LoginRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var req LoginRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (req *LoginRequest) ToLoginCommand() *command.LoginCommand {
	return &command.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	}
}
