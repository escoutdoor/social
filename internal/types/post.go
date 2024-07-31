package types

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	Text      string    `json:"text"`
	UserID    uuid.UUID `json:"user_id"`
	PhotoURL  *string   `json:"photo_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostReq struct {
	Text     string `json:"text" validate:"required,min=3"`
	PhotoURL string `json:"photo_url" validate:"omitempty,url"`
}

type UpdatePostReq struct {
	Text     string `json:"text" validate:"omitempty,min=3"`
	PhotoURL string `json:"photo_url" validate:"omitempty,url"`
}
