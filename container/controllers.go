//go:build wireinject
// +build wireinject

package container

import (
	"golearn-api-template/services"

	"github.com/dino16m/golearn-core/controller"
	"github.com/dino16m/golearn-core/event"
	"github.com/google/wire"
)

// Container is a super struct which contains all the controllers
// available in this app, it will enable service locators to function.
type ControllerContainer struct {
	AuthController    controller.AuthController
	JWTAuthController controller.JWTAuthController
}

var AuthenticatorBinding = wire.Bind(new(controller.Authenticator), new(services.AuthService))
var UserServiceBinding = wire.Bind(new(controller.UserService), new(services.UserService))

var ControllerSet = wire.NewSet(
	AuthenticatorBinding,
	UserServiceBinding,
	controller.NewAuthController,
	controller.NewJWTAuthController,
	JWTAuthServiceSet,
	wire.FieldsOf(new(ServiceContainer), "AuthService", "UserService", "JWTAuthService"),
)

func ControllerProvider(services ServiceContainer, dispatcher event.Dispatcher) ControllerContainer {
	wire.Build(ControllerSet, wire.Struct(new(ControllerContainer), "*"))
	return ControllerContainer{}
}
