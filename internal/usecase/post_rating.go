package usecase

import (
	"context"
)

func (u *Usecase) PostRating(ctx context.Context, gameID, userID string, rating int32) (bool, error) {
	return true, nil
}
