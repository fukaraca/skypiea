package model

import (
	"errors"
	"fmt"
)

var (
	ErrSessionNotFound  = NewError(1101, "session not found")
	ErrSessionNotValid  = NewError(1102, "session not valid")
	ErrIncorrectCred    = NewError(1103, "credentials are wrong")
	ErrNoSuchEmail      = NewError(1104, "there is no record for provided email")
	ErrInvalidToken     = NewError(1105, "auth token is not valid")
	ErrMissingPathParam = NewError(1106, "missing path parameter")

	ErrConversationCouldNotGet        = NewError(1201, "conversations not available")
	ErrMessagesCouldNotBeReloaded     = NewError(1202, "messages could not be reloaded")
	ErrNewMessageCouldNotBeAdded      = NewError(1203, "message could not be processed")
	ErrNewConversationCouldNotBeAdded = NewError(1203, "conversation could not be processed")
)

type Error struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Stack   []error `json:"stack"`
}

func NewError(code int, message string) *Error {
	e := &Error{
		Code:    code,
		Message: message,
		Stack:   make([]error, 0),
	}
	e.WithError(e)
	return e
}

func (e *Error) Error() string {
	if len(e.Stack) == 1 {
		return fmt.Sprintf("%s(code: %d)", e.Message, e.Code)
	}
	return fmt.Sprintf("%s:%s(code: %d)", e.Message, e.getLastError().Error(), e.Code)
}

func (e *Error) Is(target error) bool {
	return errors.Is(e, target)
}

// WithError is not wrapper
func (e *Error) WithError(err error) *Error {
	e.Stack = append(e.Stack, err)
	return e
}

func (e *Error) getLastError() error {
	return e.Stack[len(e.Stack)-1]
}
