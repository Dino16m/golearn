package mails

import (
	"fmt"
	"html/template"
	"io"
	"strings"
	txt "text/template"

	"github.com/dino16m/golearn/lib/mail"
)

// Template type accepted by the mailer
type Template interface {
	Execute(wr io.Writer, data interface{}) error
}

type data struct {
	Appname string
	Name    string
	Link    string
}

// AuthMail ...
type AuthMail struct {
	emailVerifTxt         Template
	emailVerifHTML        Template
	passwordResetCodeTxt  Template
	passwordResetCodeHTML Template
	passwordResetLinkTxt  Template
	passwordResetLinkHTML Template
	mailer                mail.IMailer
	appName               string
}

// NewAuthMail constructor of AuthMail
/*
	emailVerifTxt and emailVerifHTML templates should expect data containing fields
	{Appname string, Name string, VerificationLink string}
*/
/*
	passwordResetHTML and passwordResetTxt templates should expect data containing fields
	{Appname string, Name  string, ResetLink string or Code depending on whether it
	expects a code or a link}
*/
func NewAuthMail(appName string, mailer mail.IMailer, emailVerifTxtPath,
	emailVerifHTMLPath string, passwordResetCodeTxtPath,
	passwordResetCodeHTMLPath, passwordResetLinkTxtPath,
	passwordResetLinkHTMLPath string) (AuthMail, error) {

	emailVerifTxtTemplate, err := txt.ParseFiles(emailVerifTxtPath)
	if err != nil {
		return AuthMail{}, nil
	}
	emailVerifHTMLTemplate, err := template.ParseFiles(emailVerifHTMLPath)
	if err != nil {
		return AuthMail{}, nil
	}
	passwordResetCodeTxtTemplate, err := txt.ParseFiles(passwordResetCodeTxtPath)
	if err != nil {
		return AuthMail{}, nil
	}
	passwordResetCodeHTMLTemplate, err := template.ParseFiles(passwordResetCodeHTMLPath)
	if err != nil {
		return AuthMail{}, nil
	}
	passwordResetLinkTxtTemplate, err := txt.ParseFiles(passwordResetLinkTxtPath)
	if err != nil {
		return AuthMail{}, nil
	}
	passwordResetLinkHTMLTemplate, err := template.ParseFiles(passwordResetLinkHTMLPath)
	if err != nil {
		return AuthMail{}, nil
	}
	return AuthMail{
		appName:               appName,
		emailVerifTxt:         emailVerifTxtTemplate,
		emailVerifHTML:        emailVerifHTMLTemplate,
		passwordResetCodeTxt:  passwordResetCodeTxtTemplate,
		passwordResetCodeHTML: passwordResetCodeHTMLTemplate,
		passwordResetLinkTxt:  passwordResetLinkTxtTemplate,
		passwordResetLinkHTML: passwordResetLinkHTMLTemplate,
		mailer:                mailer,
	}, nil
}

// SendWelcomeMail ...
func (m AuthMail) SendWelcomeMail(name string, email string) {
	fmt.Println("Hi, ", name, " Welcome to ", email)
}

// SendEmailVerification sends a verification link in an email provided,
// filling the mail templates with useername and link as passed.
func (m AuthMail) SendEmailVerification(
	name string, email string, link string) error {
	data := data{m.appName, name, link}
	htmlBuilder := new(strings.Builder)
	m.emailVerifHTML.Execute(htmlBuilder, data)
	txtBuilder := new(strings.Builder)
	m.emailVerifTxt.Execute(txtBuilder, data)
	htmlMsg := htmlBuilder.String()
	txtMsg := txtBuilder.String()
	subject := fmt.Sprintf("%s Email Verification", m.appName)
	return m.sendMsg(subject, htmlMsg, txtMsg, email)
}

// SendResetLink Sends password reset link as a mail to the email provided,
// filling the mail templates with username and link as passed.
func (m AuthMail) SendResetLink(
	name string, email string, link string) error {

	data := data{m.appName, name, link}
	htmlBuilder := new(strings.Builder)
	m.passwordResetLinkHTML.Execute(htmlBuilder, data)
	txtBuilder := new(strings.Builder)
	m.passwordResetLinkTxt.Execute(txtBuilder, data)
	htmlMsg := htmlBuilder.String()
	txtMsg := htmlBuilder.String()
	subject := fmt.Sprintf("%s Password reset", m.appName)

	return m.sendMsg(subject, htmlMsg, txtMsg, email)
}

// SendResetCode Sends password reset link as a mail to the email provided,
// filling the mail templates with username and link as passed.
func (m AuthMail) SendResetCode(name, email string, code int) error {
	data := map[string]interface{}{
		"Appname": m.appName,
		"name":    name,
		"code":    code,
	}
	htmlBuilder := new(strings.Builder)
	m.passwordResetCodeHTML.Execute(htmlBuilder, data)
	txtBuilder := new(strings.Builder)
	m.passwordResetCodeTxt.Execute(txtBuilder, data)
	htmlMsg := htmlBuilder.String()
	txtMsg := htmlBuilder.String()
	subject := fmt.Sprintf("%s Password reset", m.appName)

	return m.sendMsg(subject, htmlMsg, txtMsg, email)
}

func (m AuthMail) sendMsg(
	subject string, html string, txt string, recipient string) error {
	msg := mail.InitializeMessage()
	msg.Subject = subject
	msg.HTMLMsg = html
	msg.TxtMsg = txt
	msg.Recipients = append(msg.Recipients, recipient)
	return m.mailer.Send(msg)
}
