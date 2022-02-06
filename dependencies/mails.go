// +build wireinject

package dependencies

import (
	"github.com/dino16m/golearn/lib/mail"
	"github.com/dino16m/golearn/mails"
	"github.com/google/wire"
)

// MailsContainer ...
type MailsContainer struct {
	AuthMail mails.AuthMail
}

func provideMailer() *mail.Mailer {
	mailOpts := ProvideSuperConfig().MailOptions
	return mail.NewMailer(
		mailOpts.SenderName, mailOpts.FromEmail, mailOpts.Host,
		mailOpts.Port, mailOpts.Username, mailOpts.Password,
	)
}

func provideDummyMailer() *mail.ConsoleMailer {
	return mail.NewConsoleMailer()
}
func provideIMailer() mail.IMailer {
	env := ProvideSuperConfig().Env
	var mailer mail.IMailer
	if env == "production" {
		mailer = provideMailer()
	} else {
		mailer = provideDummyMailer()
	}
	return mailer
}
func provideAuthMail() (mails.AuthMail, error) {
	mailer := provideIMailer()
	appName := ProvideSuperConfig().AppName
	t := ProvideSuperConfig().AuthMailTemplates
	return mails.NewAuthMail(appName,
		mailer, t.EmailVerifTxt,
		t.EmailVerifHTML,
		t.PasswordResetCodeTxt,
		t.PasswordResetCodeHTML, t.PasswordResetLinkTxt,
		t.PasswordResetLinkHTML)

}

// InitMails creates all the mails for the container
func InitMails() (MailsContainer, error) {
	wire.Build(
		wire.Struct(new(MailsContainer), "*"),
		provideAuthMail,
	)
	return MailsContainer{}, nil
}
