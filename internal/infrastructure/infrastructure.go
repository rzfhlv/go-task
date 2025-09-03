package infrastructure

import (
	"context"

	"github.com/rzfhlv/go-task/config"
	"github.com/rzfhlv/go-task/internal/infrastructure/sqlstore"
)

type Infrastructure interface {
	SQLStore() sqlstore.SQLStore
}

type Infra struct {
	sqlStore sqlstore.SQLStore
}

func New(ctx context.Context, cfg *config.Configuration) (Infrastructure, error) {
	sqlStore, err := sqlstore.New(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}

	return &Infra{
		sqlStore: *sqlStore,
	}, nil
}

func (i *Infra) SQLStore() sqlstore.SQLStore {
	return i.sqlStore
}
