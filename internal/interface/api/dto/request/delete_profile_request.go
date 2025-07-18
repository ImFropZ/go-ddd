package request

import (
	"encoding/json"
	"github/imfropz/go-ddd/internal/application/command"
	"io"
	"net/http"
)

type DeleteProfileRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func NewDeleteProfileRequest(r *http.Request) (*DeleteProfileRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var req DeleteProfileRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (req *DeleteProfileRequest) ToDeleteProfileCommand() *command.DeleteProfileCommand {
	return &command.DeleteProfileCommand{
		Email:    req.Email,
		Password: req.Password,
	}
}
