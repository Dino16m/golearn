package mails

import (
	"io"
)

// Template type accepted by the mailer
type Template interface {
	Execute(wr io.Writer, data interface{}) error
}
