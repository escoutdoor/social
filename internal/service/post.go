package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/escoutdoor/social/internal/cache"
	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type PostService struct {
	store store.PostStorer
	cache cache.Store
}

func NewPostService(store store.PostStorer, cache cache.Store) *PostService {
	return &PostService{
		store: store,
		cache: cache,
	}
}

func (s *PostService) Create(ctx context.Context, userID uuid.UUID, input types.CreatePostReq) (*types.Post, error) {
	post, err := s.store.Create(ctx, userID, input)
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
		p, err = s.store.GetByID(ctx, postID)
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

	post, err := s.store.Update(ctx, postID, *p)
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
		post, err = s.store.GetByID(ctx, id)
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
	key := "posts"

	posts, err := s.cache.GetPosts(ctx, key)
	if errors.Is(err, redis.Nil) {
		posts, err = s.store.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		if err := s.cache.Set(ctx, key, posts, time.Minute*1).Err(); err != nil {
			return nil, fmt.Errorf("failed to cache data: %w", err)
		}
	}
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) Delete(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	key := generatePostKey(postID)
	p, err := s.cache.GetPost(ctx, key)
	if errors.Is(err, redis.Nil) {
		p, err = s.store.GetByID(ctx, postID)
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

	if err := s.cache.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete item from cache: %w", err)
	}
	return nil
}

func generatePostKey(id uuid.UUID) string {
	return fmt.Sprintf("post%s", id)
}
