package model

import (
	"errors"
	"fmt"
)

var (
	ErrSomeTest = NewError(1001, "Some test error")
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
