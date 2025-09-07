package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	Set(ctx context.Context, key string, value int64, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) (int64, error)
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
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Del(ctx context.Context, key string) (int64, error) {
	return c.client.Del(ctx, key).Result()
}
