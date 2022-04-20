//go:build wireinject
// +build wireinject

package container

import (
	"golearn-api-template/services"

	"golearn-api-template/config"
	"golearn-api-template/repositories"

	"github.com/dino16m/golearn-core/controller"
	coreServices "github.com/dino16m/golearn-core/services"
	"github.com/google/wire"
)

type ServiceContainer struct {
	UserService    services.UserService
	AuthService    services.AuthService
	JWTAuthService coreServices.JWTAuthService
}

var ServiceContainerSet = wire.NewSet(
	services.NewAuthService,
	services.NewUserService,
	coreServices.NewJWTAuthService,
	wire.Bind(new(services.UserRepository), new(repositories.UserRepository)),
	wire.FieldsOf(new(RepositoryContainer), "UserRepository"),
	wire.FieldsOf(new(config.SuperConfig), "JwtOptions"),
)

var JWTAuthServiceSet = wire.Bind(new(controller.JWTAuthService), new(coreServices.JWTAuthService))

func ServiceProvider(repositories RepositoryContainer, cfg config.SuperConfig) ServiceContainer {
	wire.Build(ServiceContainerSet, wire.Struct(new(ServiceContainer), "*"))
	return ServiceContainer{}
}
