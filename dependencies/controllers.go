// +build wireinject

package dependencies

import (
	"github.com/dino16m/golearn/config"
	"github.com/dino16m/golearn/controllers"
	"github.com/dino16m/golearn/forms"
	"github.com/dino16m/golearn/services"
	"github.com/dino16m/golearn/types"
	"github.com/google/wire"
)

func provideAuthController(
	app App, services ServicesContainer,
	repos RepositoriesContainer,
) controllers.AuthController {
	return controllers.NewAuthController(
		app.SessionAuthMiddleware,
		services.AuthService,
		repos.UserRepo, forms.SignUpForm{}, app.EventDispatcher,
	)
}

var emailVerifControllerSet = wire.NewSet(
	wire.Bind(
		new(controllers.EmailVerificationService),
		new(services.EmailVerifService)),
	wire.FieldsOf(new(ServicesContainer), "EmailVerifService"),
	controllers.NewEmailVerifController,
)

var passwordResetControllerSet = wire.NewSet(
	wire.Bind(
		new(controllers.PasswordResetService),
		new(services.PasswordResetService)),
	wire.FieldsOf(new(ServicesContainer), "PasswordResetService"),
	wire.FieldsOf(new(ServicesContainer), "AuthService"),
	wire.Bind(new(types.Authenticator), new(services.AuthService)),
	controllers.NewPasswordResetController,
	wire.FieldsOf(new(config.SuperConfig), "AppName"),
)

var sharedDependencies = wire.NewSet(
	wire.FieldsOf(new(RepositoriesContainer), "UserRepo"),
	ProvideSuperConfig,
)

// InitRepos create the repository container
func InitControllers(
	repos RepositoriesContainer,
	services ServicesContainer, app App) ControllersContainer {
	wire.Build(
		provideAuthController,
		wire.Struct(new(ControllersContainer), "*"),
		emailVerifControllerSet,
		passwordResetControllerSet,
		sharedDependencies,
	)
	return ControllersContainer{}
}

// Container is a super struct which contains all the controllers
// available in this app, it will enable service locators to function.
type ControllersContainer struct {
	AuthController          controllers.AuthController
	EmailVerifController    controllers.EmailVerifController
	PasswordResetController controllers.PasswordResetController
}
