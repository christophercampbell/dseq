package errors

import "fmt"

// Error types
var (
	ErrInvalidConfig   = NewError("invalid configuration")
	ErrInvalidState    = NewError("invalid state")
	ErrInvalidRequest  = NewError("invalid request")
	ErrStreamError     = NewError("stream error")
	ErrDatabaseError   = NewError("database error")
	ErrValidationError = NewError("validation error")
	ErrInitialization  = NewError("initialization error")
)

// Error represents an application error
type Error struct {
	msg  string
	code string
	err  error
}

// NewError creates a new Error
func NewError(msg string) *Error {
	return &Error{
		msg:  msg,
		code: "APP_ERROR",
	}
}

// Wrap wraps an existing error
func Wrap(err error, msg string) *Error {
	return &Error{
		msg:  msg,
		code: "APP_ERROR",
		err:  err,
	}
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s: %v", e.code, e.msg, e.err)
	}
	return fmt.Sprintf("%s: %s", e.code, e.msg)
}

// Unwrap returns the wrapped error
func (e *Error) Unwrap() error {
	return e.err
}

// Is reports whether the target error is of the same type
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.code == t.code
}
