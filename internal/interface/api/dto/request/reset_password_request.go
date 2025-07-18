package request

import (
	"encoding/json"
	"github/imfropz/go-ddd/internal/application/command"
	"io"
	"net/http"
)

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func NewResetPasswordRequest(r *http.Request) (*ResetPasswordRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var req ResetPasswordRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (request *ResetPasswordRequest) ToResetPasswordCommand() *command.ResetPasswordCommand {
	return &command.ResetPasswordCommand{
		Email: request.Email,
	}
}

type ResetPasswordWithTokenRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

func NewResetPasswordWithTokenRequest(r *http.Request) (*ResetPasswordWithTokenRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var req ResetPasswordWithTokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (request *ResetPasswordWithTokenRequest) ToResetPasswordWithTokenCommand() *command.ResetPasswordWithTokenCommand {
	return &command.ResetPasswordWithTokenCommand{
		Token:       request.Token,
		NewPassword: request.NewPassword,
	}
}
