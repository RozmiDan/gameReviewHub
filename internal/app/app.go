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
	"github.com/RozmiDan/gameReviewHub/internal/usecase"

	"github.com/RozmiDan/gameReviewHub/pkg/kafka"
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

	// repo
	pg, err := postgres.New(cfg.PostgreURL.URL, postgres.MaxPoolSize(4))
	if err != nil {
		logger.Error("Cant open database", zap.Error(err))
		os.Exit(1)
	}
	defer pg.Close()

	repo := postgres_storage.New(pg, logger)

	// grpc
	ratingService, err := rating.New(context.TODO(), logger, cfg.GrpcInfo.Address, cfg.GrpcInfo.Timeout)
	if err != nil {
		os.Exit(1)
	}

	// kafka
	kafkaProducer := kafka.NewProducer(&cfg.Kafka, logger)

	// usecase
	uc := usecase.New(ratingService, repo, logger, kafkaProducer)

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




	// ids := []string{
	// 	"06bb8cf0-c346-4681-a2f1-5b90d96120b5",
	// 	"437c3a12-3008-4f5a-848e-352b6f1386da",
	// 	"c1356b34-cd5f-4358-b4e0-94bbf5527321",
	// 	"4e3bea4b-b58d-46e0-8424-af0ed96a069b",
	// 	"0905184c-0445-44b9-8a75-9cffab9a85b9",
	// 	// "8a72872d-f514-4b28-8a90-ea554ca90616",
	// 	// "920ffc7b-91c5-4480-afb4-1a22a6ce7373",
	// 	// "1a72872d-f514-4b28-8a90-ea554ca90616",
	// 	// "220ffc7b-91c5-4480-afb4-1a22a6ce7373",
	// 	// "f683e3de-ef27-470e-b909-9e5f30d9c174",
	// }

	// list, err := repo.GetGameInfo(context.TODO(), ids)
	// if err != nil {
	// 	logger.Error("errooooooor", zap.Error(err))
	// }
	// for _, it := range list {
	// 	logger.Info("response", zap.String("g", it.Genre), zap.String("id", it.ID),
	// 		zap.String("n", it.Name),
	// 	)
	// }
