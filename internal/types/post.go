package types

import "github.com/google/uuid"

type Post struct {
	ID uuid.UUID `json:"id"`
}

type CreatePostReq struct {
}
