package mapper

import (
	"github/imfropz/go-ddd/common/util"
	"github/imfropz/go-ddd/internal/application/common"
	"github/imfropz/go-ddd/internal/interface/api/dto/response"
)

func ToTokenResponse(user *common.UserResult) (*response.TokenResponse, error) {
	accessToken, err := util.GenerateAccessToken(util.AccessTokenClaims{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := util.GenerateRefreshToken(util.RefreshTokenClaims{
		Id: user.Id,
	})
	if err != nil {
		return nil, err
	}

	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
