//go:build wireinject
// +build wireinject

package container

import (
	"golearn-api-template/repositories"

	"github.com/google/wire"
	"gorm.io/gorm"
)

type RepositoryContainer struct {
	UserRepository repositories.UserRepository
}

var RepositorySet = wire.NewSet(repositories.NewUserRepository)

func RepositoryProvider(db *gorm.DB) RepositoryContainer {
	wire.Build(RepositorySet, wire.Struct(new(RepositoryContainer), "*"))
	return RepositoryContainer{}
}
