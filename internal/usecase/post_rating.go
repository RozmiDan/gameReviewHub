package usecase

import (
	"context"
	"errors"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"go.uber.org/zap"
)

func (u *Usecase) PostRating(ctx context.Context, gameID, userID string, rating int32) error {
	// 1) забираем request_id
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)

	// 2) оборачиваем логгер
	logger := u.logger.With(zap.String("func", "GetListComments"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	// проверяем наличие gameID в таблице game
	result, err := u.gameHubRepo.GetGameInfo(ctx, []string{gameID})
	if err != nil {
		if errors.Is(err, entity.ErrInternal) {
			// репозиторий упал
			logger.Error("failed to check game existence", zap.Error(err))
			return entity.ErrInternal
		}
		logger.Error("failed to check game existence", zap.Error(err))
		return entity.ErrUnidentified
	}

	if len(result) != 1 {
		// ни одной игры не вернулось
		logger.Info("game not found", zap.String("game_id", gameID))
		return entity.ErrGameNotFound
	}

	msg := entity.RatingMessage{
		GameID: gameID,
		UserID: userID,
		Rating: rating,
	}

	if err := u.kafka.PublishRating(ctx, msg); err != nil {
		logger.Error("failed to publish rating to kafka", zap.Error(err), zap.String("game_id", gameID))
		return entity.ErrBrokerUnavailable
	}

	logger.Info("rating published to kafka", zap.String("game_id", gameID))
	return nil
}
