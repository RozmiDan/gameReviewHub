package app

import (
	"context"
	"os"

	"github.com/RozmiDan/gameReviewHub/db"
	"github.com/RozmiDan/gameReviewHub/internal/config"
	"github.com/RozmiDan/gameReviewHub/pkg/logger"
	"github.com/RozmiDan/gameReviewHub/pkg/postgres"
	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) {

	logger := logger.NewLogger(cfg.Env)

	logger.Info("App started")
	logger.Debug("debug mode")

	pgxConf := pgx.ConnConfig{
		Host:     cfg.PostgreURL.Host,
		Port:     cfg.PostgreURL.Port,
		Database: cfg.PostgreURL.Database,
		User:     cfg.PostgreURL.User,
		Password: cfg.PostgreURL.Password,
	}

	db.SetupPostgres(pgxConf, logger)
	logger.Info("Migrations completed successfully\n")

	pg, err := postgres.New(cfg.PostgreURL.URL, postgres.MaxPoolSize(4))
	if err != nil {
		logger.Error("Cant open database",
			zap.Error(err),
		)
		os.Exit(1)
	}

	defer pg.Close()

	var userID int
	logger.Info("Connected postgres\n")
	query, args, err := pg.Builder.Insert("users").
		Columns("nickname", "mail").
		Values("savvy", "hfs@mail.com").
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		logger.Error("ошибка при создании запроса:",
			zap.Error(err))
	}

	err = pg.Pool.QueryRow(context.Background(), query, args...).Scan(&userID)

	if err != nil {
		logger.Info("result of db", zap.Error(err))
	}

	logger.Info("result of db", zap.Int("id", userID))
}
