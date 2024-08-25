package service

import (
	"context"

	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
)

type CommentService struct {
	store     store.CommentStorer
	postStore store.PostStorer
}

func NewCommentService(store store.CommentStorer, postStore store.PostStorer) *CommentService {
	return &CommentService{
		store:     store,
		postStore: postStore,
	}
}

func (s *CommentService) Create(ctx context.Context, userID uuid.UUID, postID uuid.UUID, input types.CreateCommentReq) (uuid.UUID, error) {
	return s.store.Create(ctx, userID, postID, input)
}

func (s *CommentService) GetByID(ctx context.Context, id uuid.UUID) (*types.Comment, error) {
	return s.store.GetByID(ctx, id)
}

func (s *CommentService) GetAll(ctx context.Context, postID uuid.UUID) ([]types.Comment, error) {
	if _, err := s.postStore.GetByID(ctx, postID); err != nil {
		return nil, err
	}
	return s.store.GetAll(ctx, postID)
}

func (s *CommentService) Delete(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	comment, err := s.store.GetByID(ctx, commentID)
	if err != nil {
		return err
	}
	if comment.UserID != userID {
		return ErrAccessDenied
	}
	return s.store.Delete(ctx, commentID)
}
