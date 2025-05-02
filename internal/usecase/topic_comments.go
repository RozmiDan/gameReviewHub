package usecase

import (
	"context"

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

	// 3) получаем метаданные игры из БД
	commentsList, err := u.gameHubRepo.GetCommentsGame(ctx, gameID, limit, offset)
	if err != nil {

		logger.Error("failed to ", zap.Error(err))
		return []entity.Comment{}, err
	}

	return commentsList, nil
}
