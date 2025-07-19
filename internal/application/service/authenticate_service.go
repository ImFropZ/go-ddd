package service

import (
	"context"
	"errors"
	"fmt"
	"github/imfropz/go-ddd/common/util"
	"github/imfropz/go-ddd/internal/application/command"
	"github/imfropz/go-ddd/internal/application/mapper"
	"github/imfropz/go-ddd/internal/domain/entity"
	"github/imfropz/go-ddd/internal/domain/event"
	"github/imfropz/go-ddd/internal/domain/repository"
	"github/imfropz/go-ddd/internal/infrastructure/db/valkey"
	"time"
)

type AuthenticateService struct {
	eventPublisher   event.EventPublisher
	valkeyRepository valkey.ValkeyRepository
	userRepository   repository.UserRepository
}

func NewAuthenticateService(eventPublisher event.EventPublisher, valkeyRepository valkey.ValkeyRepository, userRepository repository.UserRepository) *AuthenticateService {
	return &AuthenticateService{
		eventPublisher:   eventPublisher,
		valkeyRepository: valkeyRepository,
		userRepository:   userRepository,
	}
}

func (service *AuthenticateService) Profile(profileCommand *command.ProfileCommand) (*command.ProfileCommandResult, error) {
	user, err := service.userRepository.FindByEmail(profileCommand.Email)
	if err != nil {
		return nil, err
	}

	result := command.ProfileCommandResult{
		Result: mapper.NewUserResultFromEntity(user),
	}

	return &result, nil
}

func (service *AuthenticateService) Register(registerCommand *command.RegisterCommand) (*command.RegisterCommandResult, error) {
	userEntity := entity.NewUser(registerCommand.Name, registerCommand.Email, registerCommand.Password)

	validatedUser, err := entity.NewValidatedUser(userEntity)
	if err != nil {
		return nil, err
	}

	validatedUser.Password, err = util.HashPwd(validatedUser.Password)
	if err != nil {
		return nil, err
	}

	user, err := service.userRepository.Create(validatedUser)
	if err != nil {
		return nil, err
	}

	result := command.RegisterCommandResult{
		Result: mapper.NewUserResultFromEntity(user),
	}

	return &result, nil
}

func (service *AuthenticateService) Login(loginCommand *command.LoginCommand) (*command.LoginCommandResult, error) {
	user, err := service.userRepository.FindByEmail(loginCommand.Email)
	if err != nil {
		return nil, err
	}

	if err := util.ComparePwd(loginCommand.Password, user.Password); err != nil {
		return nil, err
	}

	result := command.LoginCommandResult{
		Result: mapper.NewUserResultFromEntity(user),
	}

	return &result, nil
}

func (service *AuthenticateService) UpdateProfile(updateProfileCommand *command.UpdateProfileCommand) (*command.UpdateProfileCommandResult, error) {
	old_user, err := service.userRepository.FindById(updateProfileCommand.Id)
	if err != nil {
		return nil, err
	}

	user := entity.NewUser(updateProfileCommand.Name, updateProfileCommand.Email, old_user.Password)
	user.Id = old_user.Id

	if updateProfileCommand.CurrentPassword != "" {
		if err := util.ComparePwd(updateProfileCommand.CurrentPassword, old_user.Password); err != nil {
			return nil, err
		}

		user.Password = updateProfileCommand.NewPassword
	}

	validatedUser, err := entity.NewValidatedUser(user)
	if err != nil {
		return nil, err
	}

	if updateProfileCommand.CurrentPassword != "" {
		validatedUser.Password, err = util.HashPwd(validatedUser.Password)
		if err != nil {
			return nil, err
		}
	}

	user, err = service.userRepository.Update(validatedUser)
	if err != nil {
		return nil, err
	}

	result := command.UpdateProfileCommandResult{
		Result: mapper.NewUserResultFromEntity(user),
	}

	return &result, nil
}

func (service *AuthenticateService) ResetPassword(resetPasswordCommand *command.ResetPasswordCommand) (*command.ResetPasswordCommandResult, error) {
	user, err := service.userRepository.FindByEmail(resetPasswordCommand.Email)
	if err != nil {
		return nil, err
	}

	token, err := util.GenerateResetPasswordToken(util.ResetPasswordTokenClaims{
		Email: user.Email,
	})
	if err != nil {
		return nil, err
	}

	_1_hour := 60 * 60
	service.valkeyRepository.Set(context.Background(), fmt.Sprintf("user:%s:%s", user.Id, entity.RESET_PASSWORD), token, _1_hour)

	event := entity.ResetPasswordEvent{
		Email: user.Email,
		Token: token,
		Exp:   time.Now().Add(time.Hour * time.Duration(1)),
	}

	service.eventPublisher.PublishWithKey(entity.RESET_PASSWORD, []byte(user.Email), event)

	result := command.ResetPasswordCommandResult{
		Result: mapper.NewUserResultFromEntity(user),
	}

	return &result, nil
}

func (service *AuthenticateService) ResetPasswordWithToken(resetPasswordWithTokenCommand *command.ResetPasswordWithTokenCommand) (*command.ResetPasswordWithTokenCommandResult, error) {
	claims, err := util.ValidateResetPasswordToken(resetPasswordWithTokenCommand.Token)
	if err != nil {
		return nil, err
	}

	old_user, err := service.userRepository.FindByEmail(claims.Email)
	if err != nil {
		return nil, err
	}

	cacheToken, err := service.valkeyRepository.Get(context.Background(), fmt.Sprintf("user:%s:%s", old_user.Id, entity.RESET_PASSWORD))
	if err != nil {
		return nil, err
	}

	if cacheToken != resetPasswordWithTokenCommand.Token {
		return nil, errors.New("invalid token")
	}

	user := entity.NewUser(old_user.Name, old_user.Email, resetPasswordWithTokenCommand.NewPassword)
	user.Id = old_user.Id

	validatedUser, err := entity.NewValidatedUser(user)
	if err != nil {
		return nil, err
	}

	validatedUser.Password, err = util.HashPwd(validatedUser.Password)
	if err != nil {
		return nil, err
	}

	user, err = service.userRepository.Update(validatedUser)
	if err != nil {
		return nil, err
	}

	service.valkeyRepository.Delete(context.Background(), fmt.Sprintf("user:%s:%s", old_user.Id, entity.RESET_PASSWORD))

	result := command.ResetPasswordWithTokenCommandResult{
		Result: mapper.NewUserResultFromEntity(user),
	}

	return &result, nil
}

func (service *AuthenticateService) DeleteProfile(deleteProfileCommand *command.DeleteProfileCommand) error {
	user, err := service.userRepository.FindByEmail(deleteProfileCommand.Email)
	if err != nil {
		return err
	}

	if err := util.ComparePwd(deleteProfileCommand.Password, user.Password); err != nil {
		return err
	}

	return service.userRepository.Delete(user.Id)
}
