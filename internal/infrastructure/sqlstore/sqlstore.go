package sqlstore

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/rzfhlv/go-task/config"
)

var (
	sqlDB  *sqlx.DB
	sqlErr error
	once   sync.Once
)

type SQLStore struct {
	db     *sqlx.DB
	driver string
}

func New(ctx context.Context, dbConfig config.DatabaseConfiguration) (*SQLStore, error) {
	once.Do(func() {
		var err error
		dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
			dbConfig.Driver,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Name)

		sqlDB, err = sqlx.Open(dbConfig.Driver, dsn)
		if err != nil {
			slog.ErrorContext(ctx, "error when initiating database", slog.String("error", err.Error()))
			sqlErr = err
		}

		err = sqlDB.Ping()
		if err != nil {
			slog.ErrorContext(ctx, "error when pinging to database", slog.String("error", err.Error()))
			sqlErr = err
		}
	})

	if sqlErr != nil {
		return nil, sqlErr
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
