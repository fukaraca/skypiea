package model

import (
	"errors"
	"fmt"
)

var (
	ErrSessionNotFound = NewError(1101, "session not found")
	ErrSessionNotValid = NewError(1102, "session not valid")
	ErrIncorrectCred   = NewError(1103, "credentials are wrong")
	ErrInvalidCred     = NewError(1104, "credentials are not valid")
	ErrInvalidToken    = NewError(1105, "auth token is not valid")
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func (e *Error) Is(target error) bool {
	return errors.Is(e, target)
}
