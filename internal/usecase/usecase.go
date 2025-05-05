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
	kafka        RatingProducer
}

type RatingClient interface {
	SubmitRating(ctx context.Context, userID, gameID string, rating int32) (bool, error)
	GetGameRating(ctx context.Context, gameID string) (*entity.GameRating, error)
	GetTopGames(ctx context.Context, limit, offset int32) ([]entity.GameRating, error)
}

type GameRepository interface {
	GetGameTopic(ctx context.Context, gameID string) (*entity.Game, error)
	GetGameInfo(ctx context.Context, ids []string) ([]entity.GameInList, error)
	GetCommentsGame(ctx context.Context, gameID string, limit, offset int32) ([]entity.Comment, error)
	AddComment(ctx context.Context, gameID, userID, text string) (string, error)
	AddGameTopic(ctx context.Context, gameInfo *entity.Game) (string, error)
}

type RatingProducer interface {
	PublishRating(ctx context.Context, msg entity.RatingMessage) error
}

func New(ratingClient RatingClient, gameRepo GameRepository, logger *zap.Logger, ratingProd RatingProducer) *Usecase {

	logger = logger.With(zap.String("layer", "mainUsecase"))
	return &Usecase{
		ratingClient: ratingClient,
		gameHubRepo:  gameRepo,
		logger:       logger,
		kafka:        ratingProd,
	}
}
