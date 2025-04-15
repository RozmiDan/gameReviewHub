package db

import (
	"embed"
	"os"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func SetupPostgres(conn pgx.ConnConfig, logger *zap.Logger) {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("can't set dialect in goose",
			zap.Error(err),
		)
		os.Exit(1)
	}

	db := stdlib.OpenDB(conn)
	if err := goose.Up(db, "migrations"); err != nil {
		logger.Error("can't setup migrations",
			zap.Error(err),
		)
		os.Exit(1)
	}
}
