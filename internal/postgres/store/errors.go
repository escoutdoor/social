package store

import (
	"errors"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrEmailAlreadyExists = errors.New("user with this email address already exists")
	ErrInvalidEmailOrPw   = errors.New("invalid email or password")

	ErrInvalidToken = errors.New("invalid token")

	ErrPostNotFound = errors.New("post not found")

	ErrCommentNotFound = errors.New("comment not found")

	ErrLikeFailed       = errors.New("failed to like")
	ErrRemoveLikeFailed = errors.New("failed to remove like")
	ErrAlreadyLiked     = errors.New("already liked by user")
)
