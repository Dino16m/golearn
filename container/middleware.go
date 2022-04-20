//go:build wireinject
// +build wireinject

package container

import (
	"golearn-api-template/config"
	"golearn-api-template/repositories"

	"github.com/dino16m/golearn-core/middlewares"
	"github.com/google/wire"
)

type MiddlewareContainer struct {
	CSRFMiddleware    middlewares.CSRFMiddleware
	JWTAuthMiddleware middlewares.JWTAuthMiddleware
	CORSMiddleware    middlewares.CORSMiddleware
}

func ProvideCSRFMiddleware(cfg config.SuperConfig) middlewares.CSRFMiddleware {
	return middlewares.NewCSRFMiddleware(cfg.SecretKey, cfg.Env, cfg.SessionOptions)
}

var UserRepositoryBinding = wire.Bind(new(middlewares.UserRepository), new(repositories.UserRepository))
var MiddlewareSet = wire.NewSet(
	ProvideCSRFMiddleware,
	UserRepositoryBinding,
	middlewares.NewCORSMiddleware,
	middlewares.NewJWTAuthMiddleware,
	wire.FieldsOf(new(RepositoryContainer), "UserRepository"),
	wire.FieldsOf(new(config.SuperConfig), "CORSConfig"),
	wire.FieldsOf(new(ServiceContainer), "JWTAuthService"))

func MiddlewareProvider(services ServiceContainer, repositories RepositoryContainer, cfg config.SuperConfig) MiddlewareContainer {
	wire.Build(MiddlewareSet, wire.Struct(new(MiddlewareContainer), "*"))
	return MiddlewareContainer{}
}
