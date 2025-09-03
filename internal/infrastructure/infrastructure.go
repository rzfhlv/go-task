package infrastructure

import (
	"context"

	"github.com/rzfhlv/go-task/config"
	"github.com/rzfhlv/go-task/internal/infrastructure/redis"
	"github.com/rzfhlv/go-task/internal/infrastructure/sqlstore"
)

type Infrastructure interface {
	SQLStore() *sqlstore.SQLStore
	Redis() *redis.Redis
}

type Infra struct {
	sqlStore *sqlstore.SQLStore
	redis    *redis.Redis
}

func New(ctx context.Context, cfg *config.Configuration) (Infrastructure, error) {
	sqlStore, err := sqlstore.New(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}

	redis, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		return nil, err
	}

	return &Infra{
		sqlStore: sqlStore,
		redis:    redis,
	}, nil
}

func (i *Infra) SQLStore() *sqlstore.SQLStore {
	return i.sqlStore
}

func (i *Infra) Redis() *redis.Redis {
	return i.redis
}
