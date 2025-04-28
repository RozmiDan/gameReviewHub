package usecase

import (
	"context"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
)


func (u *Usecase) GetTopicGame(ctx context.Context, gameID string) (entity.Game, error) {
	return entity.Game{}, nil
}

