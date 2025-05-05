package usecase

import (
	"context"
	"errors"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"go.uber.org/zap"
)

func (u *Usecase) CreateGameTopic(ctx context.Context, game *entity.Game) (string, error) {
	// 1) забираем request_id
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)

	// 2) оборачиваем логгер
	logger := u.logger.With(zap.String("func", "CreateGameTopic"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	gameId, err := u.gameHubRepo.AddGameTopic(ctx, game)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrGameAlreadyExists):
			logger.Info("game already exists, cannot create duplicate", zap.String("name", game.Name))
			return "", entity.ErrGameAlreadyExists

		case errors.Is(err, entity.ErrInsertGame):
			logger.Error("failed to insert game into database",
				zap.String("name", game.Name),
				zap.Error(err),
			)
			return "", entity.ErrInsertGame

		default:
			logger.Error("unexpected error creating game", zap.Error(err))
			return "", entity.ErrInternal
		}
	}

	logger.Info("game created successfully", zap.String("game_id", gameId))

	return gameId, nil
}
