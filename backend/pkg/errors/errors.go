package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrDuplicateEntry    = errors.New("duplicate entry")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInternalServer    = errors.New("internal server error")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrOrderCancelled    = errors.New("order cannot be cancelled")
	ErrInvalidStatus     = errors.New("invalid status transition")
)

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(err error, code, message string) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func Wrapf(err error, code, format string, args ...interface{}) *AppError {
	return &AppError{Code: code, Message: fmt.Sprintf(format, args...), Err: err}
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
