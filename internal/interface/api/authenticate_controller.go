package api

import (
	"encoding/json"
	"fmt"
	"github/imfropz/go-ddd/common/util"
	"github/imfropz/go-ddd/internal/application/command"
	"github/imfropz/go-ddd/internal/application/interfaces"
	"github/imfropz/go-ddd/internal/domain/repository"
	"github/imfropz/go-ddd/internal/interface/api/dto/mapper"
	"github/imfropz/go-ddd/internal/interface/api/dto/request"
	"github/imfropz/go-ddd/internal/interface/api/middleware"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type AuthenticateController struct {
	service interfaces.AuthenticateService
}

func NewAuthenticateController(r *mux.Router, service interfaces.AuthenticateService, userRepository repository.UserRepository) *AuthenticateController {
	controller := AuthenticateController{
		service: service,
	}

	r.Handle("/api/v1/profile", middleware.AuthenticationHandler(http.HandlerFunc(controller.ProfileV1), userRepository)).Methods(http.MethodGet)
	r.Handle("/api/v1/login", http.HandlerFunc(controller.LoginV1)).Methods(http.MethodPost)
	r.Handle("/api/v1/register", http.HandlerFunc(controller.RegisterV1)).Methods(http.MethodPost)
	r.Handle("/api/v1/update-profile", middleware.AuthenticationHandler(http.HandlerFunc(controller.UpdateProfileV1), userRepository)).Methods(http.MethodPost)
	r.Handle("/api/v1/reset-password", http.HandlerFunc(controller.ResetPasswordV1)).Methods(http.MethodPost)
	r.Handle("/api/v1/reset-password-with-token", http.HandlerFunc(controller.ResetPasswordWithTokenV1)).Methods(http.MethodPost)
	r.Handle("/api/v1/refresh-token", middleware.AuthenticationHandler(http.HandlerFunc(controller.RefreshTokenV1), userRepository)).Methods(http.MethodPost)
	r.Handle("/api/v1/delete-profile", middleware.AuthenticationHandler(http.HandlerFunc(controller.DeleteProfileV1), userRepository)).Methods(http.MethodPost)

	return &controller
}

func (ac *AuthenticateController) ProfileV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	claims := r.Context().Value(util.AccessTokenClaims{}).(util.AccessTokenClaims)

	profileCommand := command.ProfileCommand{
		Email: claims.Email,
	}

	user, err := ac.service.Profile(&profileCommand)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := mapper.ToUserResponse(user.Result)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ac *AuthenticateController) LoginV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	req, err := request.NewLoginRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loginCommand := req.ToLoginCommand()
	user, err := ac.service.Login(loginCommand)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := mapper.ToTokenResponse(user.Result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ac *AuthenticateController) RegisterV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	req, err := request.NewRegisterRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	registerCommand := req.ToRegisterCommand()
	user, err := ac.service.Register(registerCommand)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := mapper.ToTokenResponse(user.Result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ac *AuthenticateController) UpdateProfileV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	claims := r.Context().Value(util.AccessTokenClaims{}).(util.AccessTokenClaims)

	req, err := request.NewUpdateProfileRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	command := req.ToUpdateProfileCommand(claims.Id)
	result, err := ac.service.UpdateProfile(command)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := mapper.ToUserResponse(result.Result)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ac *AuthenticateController) ResetPasswordV1(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewResetPasswordRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	command := req.ToResetPasswordCommand()
	_, err = ac.service.ResetPassword(command)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (ac *AuthenticateController) ResetPasswordWithTokenV1(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewResetPasswordWithTokenRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	command := req.ToResetPasswordWithTokenCommand()
	_, err = ac.service.ResetPasswordWithToken(command)
	if err != nil {
		slog.Error(fmt.Sprintf("error on reset password with token: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ac *AuthenticateController) RefreshTokenV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	claims := r.Context().Value(util.AccessTokenClaims{}).(util.AccessTokenClaims)

	req, err := request.NewRefreshTokenRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if refreshClaims, err := util.ValidateRefreshToken(req.RefreshToken); err != nil && refreshClaims.Id != claims.Id {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := ac.service.Profile(&command.ProfileCommand{Email: claims.Email})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := mapper.ToTokenResponse(result.Result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ac *AuthenticateController) DeleteProfileV1(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewDeleteProfileRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deleteProfileCommand := req.ToDeleteProfileCommand()
	if err := ac.service.DeleteProfile(deleteProfileCommand); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
