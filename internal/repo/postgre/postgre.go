package postgres_storage

import (
	postgres_storage "github.com/RozmiDan/gameReviewHub/pkg/postgres"
	"go.uber.org/zap"
)

type RatingRepository struct {
	pg     *postgres_storage.Postgres
	logger *zap.Logger
}

func New(pg *postgres_storage.Postgres, logger *zap.Logger) *RatingRepository {
	logger = logger.With(zap.String("layer", "MainRepository"))
	return &RatingRepository{pg, logger}
}
