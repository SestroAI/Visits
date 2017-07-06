package errors

import "errors"

var (
	ErrUnauthorized = errors.New("User unauthorized")
	ErrWrongDataFormat = errors.New("Incompatible data passed")
	ErrBadData = errors.New("Inconsistent data found")
	ErrServerError = errors.New("We encountered a problem. Please try again")
)
