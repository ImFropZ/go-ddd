package request

import (
	"encoding/json"
	"io"
	"net/http"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func NewRefreshTokenRequest(r *http.Request) (*RefreshTokenRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var req RefreshTokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	return &req, nil
}
