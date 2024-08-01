package types

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	DOB       *time.Time `json:"date_of_birth,omitempty"`
	Bio       *string    `json:"bio,omitempty"`
	AvatarURL *string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreateUserReq struct {
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserReq struct {
	FirstName *string `json:"first_name" validate:"omitempty,min=2"`
	LastName  *string `json:"last_name" validate:"omitempty,min=2"`
	Email     *string `json:"email" validate:"omitempty,email"`
	Password  *string `json:"password" validate:"omitempty,min=6"`
	DOB       *string `json:"date_of_birth" validate:"omitempty"`
	Bio       *string `json:"bio" validate:"omitempty"`
	AvatarURL *string `json:"avatar_url" validate:"omitempty,url"`
}
