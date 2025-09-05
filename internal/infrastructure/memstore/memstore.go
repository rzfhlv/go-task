package memstore

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

type Memstore struct {
	client *redis.Client
}

func New(ctx context.Context, redisConfig config.RedisConfiguration) (*Memstore, error) {
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

	return &Memstore{
		client: redisClient,
	}, nil
}

func (r *Memstore) GetClient() *redis.Client {
	return r.client
}

func (r *Memstore) Close() error {
	if r.client != nil {
		r.client.Close()
	}

	return nil
}
