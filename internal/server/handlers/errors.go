package handlers

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidRequestBody  = errors.New("invalid request body")
	ErrForbidden           = errors.New("access denied")

	ErrNoFileReceived = errors.New("no file is received")
	ErrUnableReadFile = errors.New("unable to read the file")
	ErrUnableSaveFile = errors.New("unable to save the file")
)
