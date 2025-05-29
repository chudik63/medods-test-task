package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"medods-test-task/config"
	"medods-test-task/pkg/logger"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type DB struct {
	*sql.DB
}

const (
	driver = "postgres"
)

func New(ctx context.Context, config *config.PostgresConfig) DB {
	logs := logger.GetLoggerFromCtx(ctx)

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%d", config.User, config.Password, config.Name, config.SSLMode, config.Host, config.Port)

	db, err := sql.Open(driver, dsn)
	if err != nil {
		logs.Fatal(ctx, "can`t connect to database", zap.String("error:", err.Error()))
	}

	if err := db.Ping(); err != nil {
		logs.Fatal(ctx, "failed connecting to database", zap.String("error:", err.Error()))
	}

	logs.Debug(ctx, "database connected", zap.String("dsn", dsn))

	return DB{db}
}
