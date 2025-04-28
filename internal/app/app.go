package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RozmiDan/gameReviewHub/db"
	"github.com/RozmiDan/gameReviewHub/internal/config"
	httpserver "github.com/RozmiDan/gameReviewHub/internal/controller/http/server"
	"github.com/RozmiDan/gameReviewHub/internal/usecase"

	"github.com/RozmiDan/gameReviewHub/pkg/logger"
	"github.com/RozmiDan/gameReviewHub/pkg/postgres"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) {

	logger := logger.NewLogger(cfg.Env)

	logger.Info("App started")
	logger.Debug("debug mode")

	db.SetupPostgres(cfg, logger)
	logger.Info("Migrations completed successfully\n")

	pg, err := postgres.New(cfg.PostgreURL.URL, postgres.MaxPoolSize(4))
	if err != nil {
		logger.Error("Cant open database", zap.Error(err))
		os.Exit(1)
	}

	defer pg.Close()

	// grpc

	// ratingService, err := rating.New(context.TODO(), logger,
	// 	cfg.GrpcInfo.Address, cfg.GrpcInfo.Timeout)

	// if err != nil {
	// 	os.Exit(1)
	// }
	// logger.Info("\n\n\n")

	// go func() {
	// 	res, _ := ratingService.GetGameRating(context.TODO(), "623ee63e-b4cc-4d3b-bd6c-f5c33411fa62")
	// 	logger.Info("result", zap.Any("rate", res.AverageRating),
	// 		zap.Any("id", res.GameID),
	// 		zap.Any("count", res.RatingsCount),
	// 	)
	// }()

	// resGames, _ := ratingService.GetTopGames(context.TODO(), 10, 0)

	// for _, it := range resGames {
	// 	logger.Info("result", zap.Any("rate", it.AverageRating),
	// 		zap.Any("id", it.GameID),
	// 		zap.Any("count", it.RatingsCount),
	// 	)
	// }

	// usecase
	uc := usecase.NewUsecase()

	// server
	server := httpserver.InitServer(cfg, logger, uc)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("starting server", zap.String("port", cfg.HttpInfo.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", zap.Error(err))
			os.Exit(1)
		}
	}()

	<-stop
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Завершаем работу сервера
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	} else {
		logger.Info("Server gracefully stopped")
	}

	logger.Info("Finishing programm")
}

// var userID int
// logger.Info("Connected postgres\n")
// query, args, err := pg.Builder.Insert("users").
// 	Columns("nickname", "mail").
// 	Values("savvy", "hfs@mail.com").
// 	Suffix("RETURNING id").
// 	ToSql()

// if err != nil {
// 	logger.Error("ошибка при создании запроса:",
// 		zap.Error(err))
// }

// err = pg.Pool.QueryRow(context.Background(), query, args...).Scan(&userID)

// if err != nil {
// 	logger.Info("result of db", zap.Error(err))
// }

// logger.Info("result of db", zap.Int("id", userID))
