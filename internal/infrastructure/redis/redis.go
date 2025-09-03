package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/rzfhlv/go-task/config"
)

var (
	redisClient *redis.Client
	once        sync.Once
	redisErr    error
)

type Redis struct {
	client *redis.Client
}

func New(ctx context.Context, redisConfig config.RedisConfiguration) (*Redis, error) {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
			Password: redisConfig.Password,
			DB:       0,
		})

		err := redisClient.Ping(context.Background()).Err()
		if err != nil {
			redisErr = err
		}
	})

	if redisErr != nil {
		return nil, redisErr
	}

	return &Redis{
		client: redisClient,
	}, nil
}

func (r *Redis) GetClient() *redis.Client {
	return r.client
}

func (r *Redis) Close() error {
	if r.client != nil {
		r.client.Close()
	}

	return nil
}
