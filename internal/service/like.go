package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/escoutdoor/social/internal/cache"
	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type LikeService struct {
	store store.LikeStorer
	cache *cache.Cache
}

func NewLikeService(store store.LikeStorer, cache *cache.Cache) *LikeService {
	return &LikeService{
		store: store,
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

	err = s.LikePost(ctx, postID, userID)
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

	return s.LikeComment(ctx, commentID, userID)
}

func (s *LikeService) RemoveLikeFromPost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	err := s.RemoveLikeFromPost(ctx, postID, userID)
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
	return s.RemoveLikeFromComment(ctx, commentID, userID)
}

func (s *LikeService) isPostLiked(ctx context.Context, postID uuid.UUID) (bool, error) {
	return s.store.IsPostLiked(ctx, postID)
}

func (s *LikeService) isCommentLiked(ctx context.Context, commentID uuid.UUID) (bool, error) {
	return s.store.IsCommentLiked(ctx, commentID)
}
