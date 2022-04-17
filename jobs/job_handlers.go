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

// UserCreatedHandler handles sending welcome
// and verification mails when users are created
type UserCreatedHandler struct {
	mailer UserCreatedMail
}

type UserCreatedMail interface {
	SendEmailVerification(name string, email string, link string) error
	SendWelcomeMail(name string, email string)
}

// NewUserCreatedHandler construct the handler
func NewUserCreatedHandler(mail UserCreatedMail) UserCreatedHandler {
	return UserCreatedHandler{mailer: mail}
}

// Handle method which makes the UserCreatedHandler implement
// the JobHandler interface
func (uc UserCreatedHandler) Handle(args HandlerArgs) error {
	name := args["name"].(string)
	email := args["email"].(string)
	verificationLink := args["verificatioLink"].(string)
	uc.mailer.SendEmailVerification(name, email, verificationLink)
	uc.mailer.SendWelcomeMail(name, email)
	return nil
}
