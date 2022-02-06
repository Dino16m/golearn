package registry

import (
	"github.com/dino16m/golearn/dependencies"
	"github.com/dino16m/golearn/jobs"
)

// RegisterJobs registers handlers for all jobs available to the apps
func RegisterJobs(
	handlers dependencies.JobHandlersContainer,
	app dependencies.App) {

	worker := app.Worker
	worker.Register(jobs.UserCreatedJob, handlers.UserCreatedHandler.Handle)

}
