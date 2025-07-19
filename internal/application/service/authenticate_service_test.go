package service_test

import (
	"errors"
	"fmt"
	"github/imfropz/go-ddd/common/util"
	"github/imfropz/go-ddd/internal/application/command"
	"github/imfropz/go-ddd/internal/application/service"
	"github/imfropz/go-ddd/internal/domain/entity"
	"github/imfropz/go-ddd/internal/domain/mocks"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthenticationService_Profile(t *testing.T) {
	user := entity.NewUser("John Doe", "test@example.com", "correct-password")
	dbUser := *user
	dbUser.Password, _ = util.HashPwd(user.Password)

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		mockUserRepo.EXPECT().FindByEmail(user.Email).Return(&dbUser, nil)

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		result, err := service.Profile(&command.ProfileCommand{
			Email: user.Email,
		})

		assert.NoError(t, err)
		assert.Equal(t, result.Result.Name, user.Name)
		assert.Equal(t, result.Result.Email, user.Email)
		assert.Equal(t, result.Result.Password, dbUser.Password)
	})

	t.Run("failure: invalid email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		mockUserRepo.EXPECT().FindByEmail(user.Email).Return(nil, errors.New("user not found"))

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		_, err := service.Profile(&command.ProfileCommand{
			Email: user.Email,
		})

		assert.Error(t, err)
	})
}

func TestAuthenticationService_Register(t *testing.T) {
	user := entity.NewUser("John Doe", "test@example.com", "correct-password")

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		mockUserRepo.EXPECT().Create(gomock.Any()).Return(user, nil)

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		result, err := service.Register(&command.RegisterCommand{
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
		})

		assert.NoError(t, err)
		assert.Equal(t, result.Result.Name, user.Name)
		assert.Equal(t, result.Result.Email, user.Email)
	})

	t.Run("failure: empty fields", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		_, err := service.Register(&command.RegisterCommand{
			Name:     "",
			Email:    "",
			Password: "",
		})

		assert.Error(t, err)
	})
}

func TestAuthenticationService_Login(t *testing.T) {
	user := entity.NewUser("John Doe", "test@example.com", "correct-password")
	dbUser := *user
	dbUser.Password, _ = util.HashPwd(user.Password)

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		mockUserRepo.EXPECT().
			FindByEmail(user.Email).
			Return(&dbUser, nil)

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		_, err := service.Login(&command.LoginCommand{
			Email:    user.Email,
			Password: user.Password,
		})

		assert.NoError(t, err)
	})

	t.Run("failure: invalid password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		mockUserRepo.EXPECT().
			FindByEmail(user.Email).
			Return(&dbUser, nil)

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		_, err := service.Login(&command.LoginCommand{
			Email:    user.Email,
			Password: user.Password + "random",
		})

		assert.Error(t, err)
	})

	t.Run("failure: no user email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		mockUserRepo.EXPECT().
			FindByEmail(user.Email).
			Return(nil, errors.New("user not found"))

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		_, err := service.Login(&command.LoginCommand{
			Email:    user.Email,
			Password: user.Password,
		})

		assert.Error(t, err)
	})
}

func TestAuthenticationService_UpdateProfile(t *testing.T) {
	user := entity.NewUser("John Doe", "test@example.com", "correct-password")
	dbUser := *user
	dbUser.Password, _ = util.HashPwd(user.Password)

	t.Run("success: full", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		newUser := entity.NewUser("Jane Doe", "example@test.com", "password-correct")
		dbNewUser := *newUser
		dbNewUser.Id = user.Id
		dbNewUser.Password, _ = util.HashPwd(dbNewUser.Password)

		mockUserRepo.EXPECT().
			FindById(user.Id).
			Return(&dbUser, nil)
		mockUserRepo.EXPECT().
			Update(gomock.Any()).
			Return(&dbNewUser, nil)

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		result, err := service.UpdateProfile(&command.UpdateProfileCommand{
			Id:              user.Id,
			Name:            newUser.Name,
			Email:           newUser.Email,
			CurrentPassword: user.Password,
			NewPassword:     newUser.Password,
		})

		assert.NoError(t, err)
		assert.Equal(t, user.Id, result.Result.Id)
		assert.Equal(t, newUser.Name, result.Result.Name)
		assert.Equal(t, newUser.Email, result.Result.Email)
		assert.Equal(t, dbNewUser.Password, result.Result.Password)
	})

	t.Run("success: without updating password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		newUser := entity.NewUser("Jane Doe", "example@test.com", user.Password)
		dbNewUser := *newUser
		dbNewUser.Id = user.Id
		dbNewUser.Password = dbUser.Password

		mockUserRepo.EXPECT().
			FindById(user.Id).
			Return(&dbUser, nil)
		mockUserRepo.EXPECT().
			Update(gomock.Any()).
			Return(&dbNewUser, nil)

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		result, err := service.UpdateProfile(&command.UpdateProfileCommand{
			Id:    user.Id,
			Name:  newUser.Name,
			Email: newUser.Email,
		})

		assert.NoError(t, err)
		assert.Equal(t, newUser.Name, result.Result.Name)
		assert.Equal(t, newUser.Email, result.Result.Email)
		assert.Equal(t, dbNewUser.Password, result.Result.Password)
	})

	t.Run("failure: invalid password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		newUser := entity.NewUser("Jane Doe", "example@test.com", "password-correct")

		mockUserRepo.EXPECT().
			FindById(user.Id).
			Return(&dbUser, nil)

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		_, err := service.UpdateProfile(&command.UpdateProfileCommand{
			Id:              user.Id,
			Name:            newUser.Name,
			Email:           newUser.Email,
			CurrentPassword: user.Password + "random",
			NewPassword:     newUser.Password,
		})

		assert.Error(t, err)
	})

	t.Run("failure: empty fields", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		newUser := entity.NewUser("", "", "")

		mockUserRepo.EXPECT().
			FindById(user.Id).
			Return(&dbUser, nil)

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		_, err := service.UpdateProfile(&command.UpdateProfileCommand{
			Id:              user.Id,
			Name:            newUser.Name,
			Email:           newUser.Email,
			CurrentPassword: user.Password,
			NewPassword:     newUser.Password,
		})

		assert.Error(t, err)
	})
}

func TestAuthenticationService_ResetPassword(t *testing.T) {
	user := entity.NewUser("John Doe", "test@example.com", "correct-password")
	dbUser := *user
	dbUser.Password, _ = util.HashPwd(user.Password)

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		mockUserRepo.EXPECT().
			FindByEmail(user.Email).
			Return(&dbUser, nil)
		mockValkeyRepo.EXPECT().
			Set(gomock.Any(), fmt.Sprintf("user:%s:%s", user.Id, entity.RESET_PASSWORD), gomock.Any(), 60*60). // ttl = 1 hour
			Return(nil)
		mockEventPub.EXPECT().PublishWithKey(entity.RESET_PASSWORD, []byte(user.Email), gomock.Any()).
			Return(nil)

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		_, err := service.ResetPassword(&command.ResetPasswordCommand{
			Email: user.Email,
		})

		assert.NoError(t, err)
	})

	t.Run("failure: wrong email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		wrongEmail := "example@test.com"

		mockUserRepo.EXPECT().
			FindByEmail(wrongEmail).
			Return(nil, errors.New("user not found"))

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		_, err := service.ResetPassword(&command.ResetPasswordCommand{
			Email: wrongEmail,
		})

		assert.Error(t, err)
		assert.Contains(t, "user not found", err.Error())
	})
}

func TestAuthenticationService_ResetPasswordWithToken(t *testing.T) {
	user := entity.NewUser("John Doe", "test@example.com", "correct-password")
	dbUser := *user
	dbUser.Password, _ = util.HashPwd(user.Password)

	validToken, _ := util.GenerateResetPasswordToken(util.ResetPasswordTokenClaims{
		Email: user.Email,
	})
	invalidToken := "invalid-token"

	validPasswrd := "valid-password"

	resetPasswordTokenKey := fmt.Sprintf("user:%s:%s", user.Id, entity.RESET_PASSWORD)

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		newUser := entity.NewUser(user.Name, user.Email, validPasswrd)
		dbNewUser := *newUser
		dbNewUser.Password, _ = util.HashPwd(newUser.Password)

		mockUserRepo.EXPECT().
			FindByEmail(user.Email).
			Return(&dbUser, nil)
		mockUserRepo.EXPECT().
			Update(gomock.Any()).
			Return(&dbNewUser, nil)
		mockValkeyRepo.EXPECT().Get(gomock.Any(), resetPasswordTokenKey).Return(validToken, nil)
		mockValkeyRepo.EXPECT().Delete(gomock.Any(), resetPasswordTokenKey).Return(nil)

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		result, err := service.ResetPasswordWithToken(&command.ResetPasswordWithTokenCommand{
			Token:       validToken,
			NewPassword: validPasswrd,
		})

		assert.NoError(t, err)
		assert.Equal(t, dbNewUser.Password, result.Result.Password)
	})

	t.Run("failure: invalid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		_, err := service.ResetPasswordWithToken(&command.ResetPasswordWithTokenCommand{
			Token:       invalidToken,
			NewPassword: validPasswrd,
		})

		assert.Error(t, err)
		assert.True(t, errors.Is(err, jwt.ErrTokenMalformed))
	})
}

func TestAuthenticationService_DeleteProfile(t *testing.T) {
	user := entity.NewUser("John Doe", "test@example.com", "correct-password")
	dbUser := *user
	dbUser.Password, _ = util.HashPwd(user.Password)

	invalidPassword := "invalid-password"

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		mockUserRepo.EXPECT().FindByEmail(user.Email).Return(&dbUser, nil)
		mockUserRepo.EXPECT().Delete(user.Id).Return(nil)

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		err := service.DeleteProfile(&command.DeleteProfileCommand{
			Email:    user.Email,
			Password: user.Password,
		})

		assert.NoError(t, err)
	})

	t.Run("failure: invalid password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockValkeyRepo := mocks.NewMockValkeyRepository(ctrl)
		mockEventPub := mocks.NewMockEventPublisher(ctrl)

		mockUserRepo.EXPECT().FindByEmail(user.Email).Return(&dbUser, nil)

		service := service.NewAuthenticateService(mockEventPub, mockValkeyRepo, mockUserRepo)

		err := service.DeleteProfile(&command.DeleteProfileCommand{
			Email:    user.Email,
			Password: invalidPassword,
		})

		assert.Error(t, err)
		assert.True(t, errors.Is(err, bcrypt.ErrMismatchedHashAndPassword))
	})
}
