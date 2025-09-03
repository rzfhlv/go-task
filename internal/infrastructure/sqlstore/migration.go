package sqlstore

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rzfhlv/go-task/config"
)

const migrationDir = "internal/infrastructure/sqlstore/migration"

type Migrator struct {
	migrator *migrate.Migrate
}

func openDBWrapper(dbConfig config.DatabaseConfiguration) (database.Driver, error) {
	ctx := context.Background()
	sqlDB, err := New(ctx, dbConfig)
	if err != nil {
		slog.ErrorContext(ctx, "fail init database instance", slog.String("error", err.Error()))
		return nil, err
	}

	driver, err := postgres.WithInstance(sqlDB.db.DB, &postgres.Config{})
	if err != nil {
		slog.ErrorContext(ctx, "error while initiating postgres driver", slog.String("error", err.Error()))
		return nil, err
	}

	return driver, nil
}

func NewMigrator() (*Migrator, error) {
	cfg := config.Get()
	driver, err := openDBWrapper(cfg.Database)
	if err != nil {
		return nil, err
	}

	migrator, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationDir), cfg.Database.Driver, driver)
	if err != nil {
		return nil, err
	}

	return &Migrator{
		migrator: migrator,
	}, nil
}

func (m *Migrator) Migrate() error {
	if err := m.migrator.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}

		fmt.Println("No migration performed, current schema is the latest version.")
	}
	return nil
}

func (m *Migrator) Rollback(step int) error {
	err := m.migrator.Steps(step * -1)
	if err != nil {
		fmt.Println("failed to rollback")
		return err
	}

	fmt.Println("Rollback successfully.")

	return nil
}

func (m *Migrator) Create(name string) error {
	unixTime := time.Now().Unix()

	fileNameUp := fmt.Sprintf("%s/%d_%s.up.sql", migrationDir, unixTime, name)
	fileNameDown := fmt.Sprintf("%s/%d_%s.down.sql", migrationDir, unixTime, name)

	if _, err := os.Create(fileNameUp); err != nil {
		return err
	}

	if _, err := os.Create(fileNameDown); err != nil {
		_ = os.Remove(fileNameUp)
		return err
	}

	fmt.Println("New migration files created successfully.")

	return nil
}
