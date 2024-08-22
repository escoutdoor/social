package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/escoutdoor/social/internal/config"
	"github.com/escoutdoor/social/internal/types"
	"github.com/redis/go-redis/v9"
)

var (
	ErrUnmarshalFailed = errors.New("failed to unmarshal data from rdb")
)

type Cache struct {
	*redis.Client
}

func New(cfg *config.Config) (*Cache, error) {
	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url: %w", err)
	}

	client := redis.NewClient(opts)
	if ping, err := client.Ping(context.Background()).Result(); err != nil {
		fmt.Println(ping)
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}
	return &Cache{Client: client}, nil
}

func (c *Cache) GetPost(ctx context.Context, key string) (*types.Post, error) {
	var post types.Post
	val, err := c.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(val), &post); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrUnmarshalFailed, err)
	}
	return &post, nil
}

func (c *Cache) GetPosts(ctx context.Context, key string) ([]types.Post, error) {
	var posts []types.Post
	val, err := c.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(val), &posts); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrUnmarshalFailed, err)
	}
	return posts, nil
}
