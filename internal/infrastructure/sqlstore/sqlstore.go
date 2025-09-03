package sqlstore

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/rzfhlv/go-task/config"
)

type SQLStore struct {
	db     *sqlx.DB
	driver string
}

func New(ctx context.Context, dbConfig config.DatabaseConfiguration) (*SQLStore, error) {
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		dbConfig.Driver,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name)

	sqlDB, err := sqlx.Open(dbConfig.Driver, dsn)
	if err != nil {
		slog.ErrorContext(ctx, "error when initiating database", slog.String("error", err.Error()))
		return nil, err
	}
	return &SQLStore{
		db:     sqlDB,
		driver: dbConfig.Driver,
	}, nil
}

func (s *SQLStore) GetDB() *sqlx.DB {
	return s.db
}

func (s *SQLStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}

	return nil
}
