//go:build wireinject
// +build wireinject

package container

import (
	"golearn-api-template/config"

	"github.com/dino16m/golearn-core/event"
	gwa "github.com/gobuffalo/gocraft-work-adapter"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type App struct {
	DB              *gorm.DB
	Config          config.SuperConfig
	Controllers     ControllerContainer
	EventDispatcher *event.EventDispatcher
	Services        ServiceContainer
	Repositories    RepositoryContainer
	Middlewares     MiddlewareContainer
	Worker          *gwa.Adapter
	Logger          *logrus.Logger
}

func ProvideApp() App {
	wire.Build(
		ProvideDB, ProvideConfig, ControllerProvider,
		wire.Struct(new(App), "*"), EventSet, ServiceProvider,
		RepositoryProvider, MiddlewareProvider,
		ProvideWorker,
		ProvideLogger,
	)
	return App{}
}
