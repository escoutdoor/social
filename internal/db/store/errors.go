package store

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrEmailAlreadyExists = errors.New("user with this email address already exists")
	ErrInvalidEmailOrPw   = errors.New("invalid email or password")

	ErrInvalidToken = errors.New("invalid token")

	ErrPostNotFound = errors.New("post not found")

	ErrReplyNotFound = errors.New("reply not found")
)
