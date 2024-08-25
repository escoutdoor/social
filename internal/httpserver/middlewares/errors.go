package middlewares

import "errors"

var (
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
)
