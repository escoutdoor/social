package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/escoutdoor/social/internal/cache"
	"github.com/escoutdoor/social/internal/repository"
	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type PostService struct {
	repo  repository.Post
	cache cache.Repository
}

func NewPostService(repo repository.Post, cache cache.Repository) *PostService {
	return &PostService{
		repo:  repo,
		cache: cache,
	}
}

func (s *PostService) Create(ctx context.Context, userID uuid.UUID, input types.CreatePostReq) (*types.Post, error) {
	post, err := s.repo.Create(ctx, userID, input)
	if err != nil {
		return nil, err
	}

	key := generatePostKey(post.ID)
	if err := s.cache.Set(ctx, key, post, time.Minute*1).Err(); err != nil {
		return nil, fmt.Errorf("failed to cache data: %w", err)
	}
	return post, nil
}

func (s *PostService) Update(ctx context.Context, postID, userID uuid.UUID, input types.UpdatePostReq) (*types.Post, error) {
	key := generatePostKey(postID)
	p, err := s.cache.GetPost(ctx, key)
	if errors.Is(err, redis.Nil) {
		p, err = s.repo.GetByID(ctx, postID)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	if p.UserID != userID {
		return nil, ErrAccessDenied
	}

	if input.Content != nil {
		p.Content = *input.Content
	}
	if input.PhotoURL != nil {
		p.PhotoURL = input.PhotoURL
	}

	post, err := s.repo.Update(ctx, postID, *p)
	if err != nil {
		return nil, err
	}
	if err := s.cache.Set(ctx, key, post, time.Minute*1).Err(); err != nil {
		return nil, fmt.Errorf("failed to cache data: %w", err)
	}
	return post, nil
}

func (s *PostService) GetByID(ctx context.Context, id uuid.UUID) (*types.Post, error) {
	key := generatePostKey(id)
	post, err := s.cache.GetPost(ctx, key)
	if errors.Is(err, redis.Nil) {
		post, err = s.repo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if err := s.cache.Set(ctx, key, post, time.Minute*1).Err(); err != nil {
			return nil, fmt.Errorf("failed to cache data: %w", err)
		}
	}
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetAll(ctx context.Context) ([]types.Post, error) {
	return s.repo.GetAll(ctx)
}

func (s *PostService) Delete(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	key := generatePostKey(postID)
	p, err := s.cache.GetPost(ctx, key)
	if errors.Is(err, redis.Nil) {
		p, err = s.repo.GetByID(ctx, postID)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	if p.UserID != userID {
		return ErrAccessDenied
	}

	err = s.repo.Delete(ctx, postID)
	if err != nil {
		return err
	}

	err = s.cache.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete item from cache: %w", err)
	}
	return nil
}

func generatePostKey(id uuid.UUID) string {
	return fmt.Sprintf("post%s", id)
}
