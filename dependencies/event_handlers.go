// +build wireinject

package dependencies

import (
	"github.com/dino16m/golearn/events"
	"github.com/dino16m/golearn/mails"
	"github.com/dino16m/golearn/services"
	"github.com/google/wire"
)

type EventHandlersContainer struct {
	UserCreatedHandler events.UserCreatedHandler
}

var welcomeMailerSet = wire.NewSet(
	wire.FieldsOf(new(MailsContainer), "AuthMail"),
	wire.Bind(new(events.WelcomeMailer), new(mails.AuthMail)),
)
var verificationServiceSet = wire.NewSet(
	wire.FieldsOf(new(ServicesContainer), "EmailVerifService"),
	wire.Bind(new(events.VerificationService), new(services.EmailVerifService)),
)
var userCreatedEventHandlerSet = wire.NewSet(
	verificationServiceSet,
	welcomeMailerSet,
	events.NewUserCreatedHandler,
)

func InitEventHandlers(s ServicesContainer,
	m MailsContainer) EventHandlersContainer {
	wire.Build(
		wire.Struct(new(EventHandlersContainer), "*"),
		userCreatedEventHandlerSet,
	)
	return EventHandlersContainer{}
}
