package usecase

import (
	"context"
	"errors"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"go.uber.org/zap"
)

func (u *Usecase) GetTopicGame(ctx context.Context, gameID string) (*entity.Game, error) {
	// 1) забираем request_id
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)

	// 2) оборачиваем логгер
	logger := u.logger.With(zap.String("func", "GetTopicGame"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	// 3) получаем метаданные игры из БД
	game, err := u.gameHubRepo.GetGameTopic(ctx, gameID)
	if err != nil {
		if errors.Is(err, entity.ErrGameNotFound) {
			logger.Info("game not found in repository", zap.String("game_id", gameID))
			return &entity.Game{}, entity.ErrGameNotFound
		}
		logger.Error("failed to fetch game metadata", zap.Error(err))
		return &entity.Game{}, err
	}

	// 4) получаем рейтинг через RPC
	rating, err := u.ratingClient.GetGameRating(ctx, gameID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrGameNotFound):
			// просто нет оценок — оставляем Game.Rating = zero, продолжаем
			logger.Info("no ratings yet for game", zap.String("game_id", gameID))

		case errors.Is(err, entity.ErrInvalidUUID):
			logger.Warn("invalid gameID passed to rating service", zap.String("game_id", gameID))
			return nil, entity.ErrInvalidUUID

		case errors.Is(err, entity.ErrServiceUnavailable):
			// вернем Game без рейтинга
			logger.Error("rating service unavailable", zap.String("game_id", gameID))

		case errors.Is(err, entity.ErrInternalRating):
			logger.Error("rating service error", zap.String("game_id", gameID))

		default:
			// неожиданный сбой (например, ctx.Err() или сетевые), пробрасываем
			logger.Error("unexpected error from rating client", zap.Error(err))
			return nil, err
		}
	} else {
		// Всё ок — присоединяем рейтинг
		game.Rating = *rating
	}

	// 5) возвращаем финальный Game
	return game, nil
}
