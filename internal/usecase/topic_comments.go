package usecase

import (
	"context"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
)

func (u *Usecase) GetListComments(ctx context.Context, gameID string, limit, offset int) ([]entity.Comment, error) {
	return []entity.Comment{}, nil
}
