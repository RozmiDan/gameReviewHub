package usecase

import (
	"context"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
)

// ListGames получает топ-N игр с учётом пагинации
// func (u *Usecase) ListGames(ctx context.Context, limit, offset int) ([]entity.GameInList, error) {
// 	//(RPC → БД → merge)
// 	u.ratingClient.GetTopGames(ctx, int32(limit), int32(offset))
// }

// ListGames получает топ-N игр с учётом пагинации
func (u *Usecase) GetListGames(ctx context.Context, limit, offset int32) ([]entity.GameInList, error) {
	//(RPC → БД → merge)

	return []entity.GameInList{}, nil
}


