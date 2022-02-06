// +build wireinject

package dependencies

import (
	"github.com/dino16m/golearn/jobs"
	"github.com/dino16m/golearn/mails"
	"github.com/google/wire"
)

var userCreatedJobHandlerSet = wire.NewSet(
	wire.Bind(new(jobs.UserCreatedMail), new(mails.AuthMail)),
	wire.FieldsOf(new(MailsContainer), "AuthMail"),
	jobs.NewUserCreatedHandler,
)

// JobHandlersContainer is a super struct which contains all the repositories
// available in a repository package, it will enable service locators to function.
type JobHandlersContainer struct {
	UserCreatedHandler jobs.UserCreatedHandler
}

// InitJobs Create the jobs container
func InitJobs(mails MailsContainer) JobHandlersContainer {
	wire.Build(
		userCreatedJobHandlerSet,
		wire.Struct(new(JobHandlersContainer), "*"),
	)
	return JobHandlersContainer{}
}
