package usecase

import (
	"context"
	"errors"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"go.uber.org/zap"
)

func (u *Usecase) AddComment(ctx context.Context, gameID, userID, text string) (string, error) {
	// 1) забираем request_id
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)

	// 2) оборачиваем логгер
	logger := u.logger.With(zap.String("func", "AddComment"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	commId, err := u.gameHubRepo.AddComment(ctx, gameID, userID, text)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrGameNotFound):
			logger.Info("game not found, cannot add comment", zap.String("game_id", gameID))
			return "", entity.ErrGameNotFound

		case errors.Is(err, entity.ErrInsertComment):
			logger.Error("failed to insert comment into database",
				zap.String("game_id", gameID),
				zap.String("user_id", userID),
				zap.Error(err),
			)
			return "", entity.ErrInsertComment

		default:
			logger.Error("cannot add comment: unexpected error", zap.Error(err))
			return "", entity.ErrInternal
		}
	}

	logger.Info("comment added successfully", zap.String("comment_id", commId))

	return commId, nil
}
