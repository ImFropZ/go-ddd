package interfaces

import "github/imfropz/go-ddd/internal/application/command"

type AuthenticateService interface {
	Profile(profileCommand *command.ProfileCommand) (*command.ProfileCommandResult, error)
	Register(registerCommand *command.RegisterCommand) (*command.RegisterCommandResult, error)
	Login(loginCommand *command.LoginCommand) (*command.LoginCommandResult, error)
	ResetPassword(resetPasswordCommand *command.ResetPasswordCommand) (*command.ResetPasswordCommandResult, error)
	ResetPasswordWithToken(resetPasswordWithTokenCommand *command.ResetPasswordWithTokenCommand) (*command.ResetPasswordWithTokenCommandResult, error)
	DeleteProfile(deleteProfileCommand *command.DeleteProfileCommand) error
}
