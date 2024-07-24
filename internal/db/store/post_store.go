package store

import (
	"context"
	"database/sql"

	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
)

type PostStore struct {
	db *sql.DB
}

type PostStorer interface {
	Create(ctx context.Context, userID uuid.UUID, input types.CreatePostReq) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*types.Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewPostStore(db *sql.DB) *PostStore {
	return &PostStore{
		db: db,
	}
}

func (s *PostStore) Create(ctx context.Context, userID uuid.UUID, input types.CreatePostReq) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (s *PostStore) GetByID(ctx context.Context, id uuid.UUID) (*types.Post, error) {
	return nil, nil
}

func (s *PostStore) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
