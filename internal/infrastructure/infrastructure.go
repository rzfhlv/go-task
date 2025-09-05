package infrastructure

import (
	"context"

	"github.com/rzfhlv/go-task/config"
	"github.com/rzfhlv/go-task/internal/infrastructure/memstore"
	"github.com/rzfhlv/go-task/internal/infrastructure/sqlstore"
)

type Infrastructure interface {
	SQLStore() *sqlstore.SQLStore
	MemStore() *memstore.Memstore
}

type Infra struct {
	sqlStore *sqlstore.SQLStore
	memStore *memstore.Memstore
}

func New(ctx context.Context, cfg *config.Configuration) (Infrastructure, error) {
	sqlStore, err := sqlstore.New(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}

	memStore, err := memstore.New(ctx, cfg.Redis)
	if err != nil {
		return nil, err
	}

	return &Infra{
		sqlStore: sqlStore,
		memStore: memStore,
	}, nil
}

func (i *Infra) SQLStore() *sqlstore.SQLStore {
	return i.sqlStore
}

func (i *Infra) MemStore() *memstore.Memstore {
	return i.memStore
}
