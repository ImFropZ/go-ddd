package rest

import (
	"encoding/json"
	"github/imfropz/go-ddd/internal/application/interfaces"
	"github/imfropz/go-ddd/internal/interface/api/rest/dto/filter"
	"github/imfropz/go-ddd/internal/interface/api/rest/dto/mapper"
	"github/imfropz/go-ddd/internal/interface/api/rest/dto/request"
	"net/http"

	"github.com/gorilla/mux"
)

type UserController struct {
	service interfaces.UserService
}

func NewUserController(r *mux.Router, service interfaces.UserService) *UserController {
	controller := &UserController{
		service: service,
	}

	r.Handle("/api/v1/users", http.HandlerFunc(controller.QueryUsersControllerV1)).Methods(http.MethodGet)
	r.Handle("/api/v1/users", http.HandlerFunc(controller.CreateUserControllerV1)).Methods(http.MethodPost)

	return controller
}

func (controller *UserController) QueryUsersControllerV1(w http.ResponseWriter, r *http.Request) {
	userCriteria, err := filter.RequestToUserCriteria(*r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := controller.service.FindAllUsers(userCriteria)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := mapper.ToUserListResponse(result.Result)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (controller *UserController) CreateUserControllerV1(w http.ResponseWriter, r *http.Request) {
	req, err := request.NewCreateUserRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := controller.service.CreateUser(req.ToCreateUserCommand())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := mapper.ToUserResponse(result.Result)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
