// +build wireinject

package dependencies

import (
	"github.com/dino16m/golearn/config"
	"github.com/dino16m/golearn/services"
	"github.com/google/wire"
)

// ServicesContainer is a super struct which contains all the services
// available in this app, it will enable service locators to function.
type ServicesContainer struct {
	AuthService          services.AuthService
	EmailVerifService    services.EmailVerifService
	PasswordResetService services.PasswordResetService
}

var authServiceSet = wire.NewSet(
	wire.FieldsOf(new(RepositoriesContainer), "UserRepo"),
	services.NewAuthService,
)

func providePasswordResetService(
	mails MailsContainer,
	cfg config.SuperConfig) (services.PasswordResetService, error) {
	return services.NewPasswordResetService(
		"verify-password-reset-link",
		cfg.SecretKey, mails.AuthMail)
}

func provideEmailVerifService(mails MailsContainer,
	cfg config.SuperConfig) (services.EmailVerifService, error) {
	return services.NewEmailVerifService("verify-email",
		cfg.SecretKey, mails.AuthMail)
}

// InitServices wires up the ServicesContainer
func InitServices(repos RepositoriesContainer,
	mails MailsContainer) (ServicesContainer, error) {
	wire.Build(
		ProvideSuperConfig,
		provideEmailVerifService,
		authServiceSet,
		providePasswordResetService,
		wire.Struct(new(ServicesContainer), "*"),
	)
	return ServicesContainer{}, nil
}
