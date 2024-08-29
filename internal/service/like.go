package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/escoutdoor/social/internal/cache"
	"github.com/escoutdoor/social/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type LikeService struct {
	repo  repository.Like
	cache cache.Repository
}

func NewLikeService(repo repository.Like, cache cache.Repository) *LikeService {
	return &LikeService{
		repo:  repo,
		cache: cache,
	}
}

func (s *LikeService) LikePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	il, err := s.isPostLiked(ctx, postID)
	if err != nil {
		return err
	}
	if il {
		return ErrAlreadyLiked
	}

	err = s.repo.LikePost(ctx, postID, userID)
	if err != nil {
		return err
	}

	key := generatePostKey(postID)
	p, err := s.cache.GetPost(ctx, key)
	if errors.Is(err, redis.Nil) {
		return nil
	}
	if err != nil {
		return err
	}

	p.Likes++
	if err := s.cache.Set(ctx, generatePostKey(p.ID), p, time.Minute*1).Err(); err != nil {
		return fmt.Errorf("failed to cache data: %w", err)
	}
	return nil
}

func (s *LikeService) LikeComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	il, err := s.isCommentLiked(ctx, commentID)
	if err != nil {
		return err
	}
	if il {
		return ErrAlreadyLiked
	}

	return s.repo.LikeComment(ctx, commentID, userID)
}

func (s *LikeService) RemoveLikeFromPost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	err := s.repo.RemoveLikeFromPost(ctx, postID, userID)
	if err != nil {
		return err
	}

	key := generatePostKey(postID)
	p, err := s.cache.GetPost(ctx, key)
	if errors.Is(err, redis.Nil) {
		return nil
	}
	if err != nil {
		return err
	}

	p.Likes--
	if err := s.cache.Set(ctx, generatePostKey(p.ID), p, time.Minute*1).Err(); err != nil {
		return fmt.Errorf("failed to cache data: %w", err)
	}
	return nil
}

func (s *LikeService) RemoveLikeFromComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	return s.repo.RemoveLikeFromComment(ctx, commentID, userID)
}

func (s *LikeService) isPostLiked(ctx context.Context, postID uuid.UUID) (bool, error) {
	return s.repo.IsPostLiked(ctx, postID)
}

func (s *LikeService) isCommentLiked(ctx context.Context, commentID uuid.UUID) (bool, error) {
	return s.repo.IsCommentLiked(ctx, commentID)
}
