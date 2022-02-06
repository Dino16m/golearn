package errors

import "net/http"

type AppError struct {
	Code    int
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

var (
	ErrMissingLoginValues = AppError{
		http.StatusUnprocessableEntity, "Login values absent",
	}
	ErrAuthenticationFailed = AppError{
		http.StatusUnauthorized, "Authentication failed",
	}
)
