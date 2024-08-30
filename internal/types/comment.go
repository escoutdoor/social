package types

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID              uuid.UUID  `json:"id"`
	Content         string     `json:"content"`
	UserID          uuid.UUID  `json:"user_id"`
	PostID          uuid.UUID  `json:"post_id"`
	ParentCommentID *uuid.UUID `json:"parent_comment_id"`
	Replies         []Comment  `json:"replies,omitempty"`
	Likes           int        `json:"likes"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type CreateCommentReq struct {
	Content         string     `json:"content" validate:"required,min=3"`
	ParentCommentID *uuid.UUID `json:"parent_comment_id" validate:"omitempty,uuid"`
}
