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
	rating "github.com/RozmiDan/gameReviewHub/internal/repo/grpcclient"
	postgres_storage "github.com/RozmiDan/gameReviewHub/internal/repo/postgre"
	redis_build "github.com/RozmiDan/gameReviewHub/internal/repo/redis"
	"github.com/RozmiDan/gameReviewHub/internal/usecase"

	"github.com/RozmiDan/gameReviewHub/pkg/kafka"
	"github.com/RozmiDan/gameReviewHub/pkg/logger"
	prom_metrics "github.com/RozmiDan/gameReviewHub/pkg/metrics"
	"github.com/RozmiDan/gameReviewHub/pkg/postgres"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) {

	logger := logger.NewLogger(cfg.Env)

	logger.Info("App started")
	logger.Debug("debug mode")

	db.SetupPostgres(cfg, logger)
	logger.Info("Migrations completed successfully\n")

	// repo
	pg, err := postgres.New(cfg.PostgreURL.URL, postgres.MaxPoolSize(25))
	if err != nil {
		logger.Error("Cant open database", zap.Error(err))
		os.Exit(1)
	}
	defer pg.Close()

	repo := postgres_storage.New(pg, logger)

	// grpc
	ratingService, err := rating.New(context.TODO(), logger, cfg.GrpcInfo.Address, cfg.GrpcInfo.Timeout)
	if err != nil {
		logger.Error("Cant connect to rating service")
		os.Exit(1)
	}

	// инициализируем метрики
	prom_metrics.Init()

	// kafka
	kafkaProducer := kafka.NewProducer(&cfg.Kafka, logger)
	//kafkaProducer := &RatingProducer{}

	// redis
	redisClient := redis_build.NewRedisClient(cfg.Redis.RedisAddress, cfg.Redis.RedisPassword,
		cfg.Redis.RedisDB, cfg.Redis.RedisTTL, logger)

	// usecase
	uc := usecase.New(ratingService, repo, logger, kafkaProducer, redisClient)

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

// type RatingProducer struct {
// }

// func (*RatingProducer) PublishRating(ctx context.Context, msg entity.RatingMessage) error {
// 	return nil
// }
