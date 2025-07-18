package mapper

import (
	"github/imfropz/go-ddd/internal/application/common"
	"github/imfropz/go-ddd/internal/interface/api/dto/response"
)

func ToUserResponse(user *common.UserResult) *response.UserResponse {
	return &response.UserResponse{
		Id:        user.Id.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserListResponse(users []*common.UserResult) *response.ListUsersResponse {
	res := response.ListUsersResponse{
		Users: make([]*response.UserResponse, 0),
	}
	for _, user := range users {
		res.Users = append(res.Users, ToUserResponse(user))
	}
	return &res
}
