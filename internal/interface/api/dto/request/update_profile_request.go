package request

import (
	"encoding/json"
	"github/imfropz/go-ddd/internal/application/command"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type UpdateProfileRequest struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func NewUpdateProfileRequest(r *http.Request) (*UpdateProfileRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var req UpdateProfileRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (req *UpdateProfileRequest) ToUpdateProfileCommand(id uuid.UUID) *command.UpdateProfileCommand {
	return &command.UpdateProfileCommand{
		Id:              id,
		Name:            req.Name,
		Email:           req.Email,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}
}
