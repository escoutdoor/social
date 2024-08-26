package service

import (
	"context"

	"github.com/escoutdoor/social/internal/repository"
	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
)

type CommentService struct {
	repo     repository.Comment
	postRepo repository.Post
}

func NewCommentService(repo repository.Comment, postRepo repository.Post) *CommentService {
	return &CommentService{
		repo:     repo,
		postRepo: postRepo,
	}
}

func (s *CommentService) Create(ctx context.Context, userID uuid.UUID, postID uuid.UUID, input types.CreateCommentReq) (uuid.UUID, error) {
	return s.repo.Create(ctx, userID, postID, input)
}

func (s *CommentService) GetByID(ctx context.Context, id uuid.UUID) (*types.Comment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CommentService) GetAll(ctx context.Context, postID uuid.UUID) ([]types.Comment, error) {
	if _, err := s.postRepo.GetByID(ctx, postID); err != nil {
		return nil, err
	}
	return s.repo.GetAll(ctx, postID)
}

func (s *CommentService) Delete(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	comment, err := s.repo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}
	if comment.UserID != userID {
		return ErrAccessDenied
	}
	return s.repo.Delete(ctx, commentID)
}
