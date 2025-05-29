package migrator

import (
	"errors"
	"fmt"
	"medods-test-task/config"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Start(cfg *config.Config) error {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Name, cfg.Postgres.SSLMode)

	m, err := migrate.New("file://"+cfg.Server.MigrationsPath, dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to make migration up: %w", err)
	}

	return nil
}
