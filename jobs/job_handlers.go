package jobs

import (
	"github.com/gobuffalo/buffalo/worker"
)

// JobHandler interface expected of any struct that is
// registered as a handler for a job
type JobHandler interface {
	Handle(worker.Args) error
}

type HandlerArgs = worker.Args
