package service

import "errors"

var (
	ErrInvalidEmailOrPw = errors.New("invalid email or password")

	ErrInvalidToken = errors.New("invalid token")

	ErrAccessDenied = errors.New("access denied")

	ErrAlreadyLiked = errors.New("already liked by user")
)
