package usecase

import (
	"context"
	"errors"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"go.uber.org/zap"
)

func (u *Usecase) GetListComments(ctx context.Context, gameID string, limit, offset int32) ([]entity.Comment, error) {
	// 1) забираем request_id
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)

	// 2) оборачиваем логгер
	logger := u.logger.With(zap.String("func", "GetListComments"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	// 3) получаем комментарии
	commentsList, err := u.gameHubRepo.GetCommentsGame(ctx, gameID, limit, offset)
	if err != nil {
		// timeout
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			logger.Error("timeout fetching comments", zap.Error(err))
			return nil, entity.ErrTimeout
		}
		// внутренняя ошибка чтения комментариев
		if errors.Is(err, entity.ErrInternalComments) {
			logger.Error("internal error fetching comments", zap.Error(err))
			return nil, entity.ErrInternal
		}
		// здесь больше нечего ловить — прокидываем дальше
		logger.Error("unexpected error fetching comments", zap.Error(err))
		return nil, entity.ErrInternal
	}

	return commentsList, nil
}
