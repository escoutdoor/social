package store

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidEmailOrPw  = errors.New("invalid email or password")

	ErrPostNotFound = errors.New("post not found")
)
