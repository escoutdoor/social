package handlers

import "errors"

var (
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidRequestBody = errors.New("invalid request body")

	ErrFileNotReceived = errors.New("no file received")
	ErrFileReadFailed  = errors.New("failed to read the file")
	ErrFileSaveFailed  = errors.New("failed to save the file")
)
