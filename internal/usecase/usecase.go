package usecase

import (
	"context"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"go.uber.org/zap"
)

type Usecase struct {
	ratingClient RatingClient
	gameHubRepo  GameRepository
	logger       *zap.Logger
}

type RatingClient interface {
	SubmitRating(ctx context.Context, userID, gameID string, rating int32) (bool, error)
	GetGameRating(ctx context.Context, gameID string) (*entity.GameRating, error)
	GetTopGames(ctx context.Context, limit, offset int32) ([]entity.GameRating, error)
}

type GameRepository interface {
}

func New(ratingClient RatingClient, gameRepo GameRepository, logger *zap.Logger) *Usecase {

	logger = logger.With(zap.String("layer", "mainService"))
	return &Usecase{
		ratingClient: ratingClient,
		gameHubRepo:  gameRepo,
		logger:       logger,
	}
}
