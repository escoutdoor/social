package types

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	UserID    uuid.UUID `json:"user_id"`
	PhotoURL  *string   `json:"photo_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostReq struct {
	Content  string `json:"content" validate:"required,min=3"`
	PhotoURL string `json:"photo_url" validate:"omitempty,url"`
}

type UpdatePostReq struct {
	Content  *string `json:"content" validate:"omitempty,min=3"`
	PhotoURL *string `json:"photo_url" validate:"omitempty,url"`
}
