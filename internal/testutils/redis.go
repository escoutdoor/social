package testutils

import (
	"context"
	"fmt"

	"github.com/escoutdoor/social/internal/cache"
	"github.com/testcontainers/testcontainers-go"
	redismodule "github.com/testcontainers/testcontainers-go/modules/redis"
)

func NewRedisContainer() (testcontainers.Container, cache.Repository, error) {
	ctx := context.Background()
	container, err := redismodule.Run(ctx, "redis:6")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to run redis container: %w", err)
	}

	p, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get container external port: %w", err)
	}

	redisAddr := fmt.Sprintf("redis://localhost:%s", p.Port())
	c, err := cache.New(redisAddr)
	if err != nil {
		return nil, nil, err
	}

	return container, c, nil
}
