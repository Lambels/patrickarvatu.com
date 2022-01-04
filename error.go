package pa

import (
	"errors"
	"fmt"
)

// Error codes which map good to http errors.
const (
	ECONFLICT       = "conflict"
	EINTERNAL       = "internal"
	EINVALID        = "invalid"
	ENOTFOUND       = "not_found"
	ENOTIMPLEMENTED = "not_implemented"
	EUNAUTHORIZED   = "unauthorized"
)

// Error is a struct containing full details about the error.
type Error struct {
	// Code to check the type of the error.
	Code string

	// Human readeable message.
	Message string
}

// Error is used to implement the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("wtf error: code=%s message=%s", e.Code, e.Message)
}

// ErrorCode is a helper function to retrieve the error code from a pa.Error.
// returns an empty string if err is nil.
// returns EINTERNAL if the error isnt a pa.Error.
func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}
	return EINTERNAL
}

// ErrorMessage is a helper function to retrieve the error message from pa.Error.
// returns an empty string if err is nil.
// returns "Internal error." if the error isnt a pa.Error.
func ErrorMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}
	return "Internal error."
}

// Errorf is a helper function to quickly init an error with code and format: message.
func Errorf(code string, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
