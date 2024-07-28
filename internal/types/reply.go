package types

import (
	"time"

	"github.com/google/uuid"
)

type Reply struct {
	ID        uuid.UUID `json:"id"`
	Text      string    `json:"text"`
	UserID    uuid.UUID `json:"user_id"`
	PostID    uuid.UUID `json:"post_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateReplyReq struct {
	Text string `json:"text" validate:"required,min=3"`
}
