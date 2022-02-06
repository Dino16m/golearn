// +build wireinject

package dependencies

import (
	"github.com/dino16m/golearn/types"
	"github.com/google/wire"
)

// RepositoriesContainer is a super struct which contains all the repositories
// available in a repository package, it will enable service locators to function.
type RepositoriesContainer struct {
	UserRepo types.UserRepository
}

// InitRepos create the repository container
func InitRepos(repo types.UserRepository) RepositoriesContainer {
	wire.Build(
		wire.Struct(new(RepositoriesContainer), "*"),
	)
	return RepositoriesContainer{}
}
