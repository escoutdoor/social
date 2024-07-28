package handlers

import "errors"

var (
	ErrIntervalServerError = errors.New("internal server error")
	ErrInvalidRequestBody  = errors.New("invalid request body")
	ErrForbidden           = errors.New("access denied")
)
