package errors

import "github.com/dino16m/golearn-core/errors"

func ValidationError(message string) errors.AppError {
	return errors.AppError{Code: 400, Message: message}
}

func Forbidden(message string) errors.AppError {
	return errors.AppError{Code: 403, Message: message}
}

func NotFound(message string) errors.AppError {
	return errors.AppError{Code: 404, Message: message}
}
