package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	Set(ctx context.Context, key string, value int64, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type Cache struct {
	client *redis.Client
}

func New(client *redis.Client) CacheRepository {
	return &Cache{
		client: client,
	}
}

func (c *Cache) Set(ctx context.Context, key string, value int64, duration time.Duration) error {
	return c.client.Set(ctx, key, value, duration).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	stmd := c.client.Get(ctx, key)
	return stmd.Val(), stmd.Err()
}
