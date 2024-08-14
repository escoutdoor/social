package types

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	UserID    uuid.UUID `json:"user_id"`
	PhotoURL  *string   `json:"photo_url,omitempty"`
	Likes     int       `json:"likes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p Post) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Post) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

type Posts []Post

func (p Posts) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Posts) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

type CreatePostReq struct {
	Content  string `json:"content" validate:"required,min=3"`
	PhotoURL string `json:"photo_url" validate:"omitempty,url"`
}

type UpdatePostReq struct {
	Content  *string `json:"content" validate:"omitempty,min=3"`
	PhotoURL *string `json:"photo_url" validate:"omitempty,url"`
}
