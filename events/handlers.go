package events

import (
	"github.com/dino16m/golearn/lib/event"
)

// UserCreatedHandler  ...
type UserCreatedHandler struct {
	*event.BaseHandler
	mail     WelcomeMailer
	verifier VerificationService
}

// VerificationService interface expected by the listener to make url
// to verify a user.
type VerificationService interface {
	SendVerificationMail(name string, id int, email, baseURL string) error
}

// WelcomeMailer ...
type WelcomeMailer interface {
	SendWelcomeMail(name string, email string)
}

// NewUserCreatedHandler constructs the UserCreatedHandler struct
func NewUserCreatedHandler(mail WelcomeMailer,
	verifService VerificationService) UserCreatedHandler {
	return UserCreatedHandler{
		mail: mail, verifier: verifService, BaseHandler: &event.BaseHandler{},
	}
}

// Handle is called when the event listeners are executed
func (uch UserCreatedHandler) Handle(args interface{}) {
	params := args.(map[string]interface{})
	firstName := params["firstName"].(string)
	email := params["email"].(string)
	id := params["id"].(int)
	baseURL := params["baseURL"].(string)

	uch.verifier.SendVerificationMail(firstName, id, email, baseURL)
	uch.mail.SendWelcomeMail(firstName, email)
}
