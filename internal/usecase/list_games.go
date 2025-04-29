package usecase

import (
	"context"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"go.uber.org/zap"
)

// ListGames получает топ-N игр с учётом пагинации
func (u *Usecase) GetListGames(ctx context.Context, limit, offset int32) ([]entity.GameInList, error) {
	//(RPC → БД → merge)

	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)
	logger := u.logger
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}
	logger = logger.With(zap.String("func", "GetListGames"))

	// RPC
	ratings, err := u.ratingClient.GetTopGames(ctx, limit, offset)
	if err != nil {
		logger.Error("failed to fetch top games from rating service", zap.Error(err))
		return nil, err
	}
	if len(ratings) == 0 {
		return []entity.GameInList{}, nil
	}

	// DB
	ids := make([]string, len(ratings))
	for i, r := range ratings {
		ids[i] = r.GameID
	}
	metas, err := u.gameHubRepo.GetGameInfo(ctx, ids)
	if err != nil {
		logger.Error("failed to fetch game metadata", zap.Error(err))
		return nil, err
	}

	metaMap := make(map[string]entity.GameInList, len(metas))
	for _, m := range metas {
		metaMap[m.ID] = m
	}

	out := make([]entity.GameInList, 0, len(ratings))
	for _, r := range ratings {
		meta, ok := metaMap[r.GameID]
		if !ok {
			//logger.Info("metadata missing for game", zap.String("game_id", r.GameID))
			continue
		}
		meta.Rating = r.AverageRating
		out = append(out, meta)
	}

	logger.Info("completed", zap.Int("returned", len(out)))
	return out, nil
}
