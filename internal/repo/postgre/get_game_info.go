package postgres_storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"go.uber.org/zap"
)

func (r *RatingRepository) GetGameInfo(ctx context.Context, ids []string) ([]entity.GameInList, error) {
	// 1) забираем request_id
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)

	// 2) оборачиваем логгер
	logger := r.logger.With(zap.String("func", "GetGameInfo"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	if len(ids) == 0 {
		logger.Info("no ids provided, returning empty list")
		return []entity.GameInList{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	q := fmt.Sprintf(`
        SELECT id, name, genre
        FROM games
        WHERE id IN (%s)
    `, strings.Join(placeholders, ","))

	rows, err := r.pg.Pool.Query(ctx, q, args...)
	if err != nil {
		logger.Error("query failed", zap.Error(err))
		return nil, entity.ErrInternal
	}
	defer rows.Close()

	var out []entity.GameInList
	for rows.Next() {
		var g entity.GameInList
		if err := rows.Scan(&g.ID, &g.Name, &g.Genre); err != nil {
			logger.Error("scan failed", zap.Error(err))
			return nil, entity.ErrInternal
		}
		out = append(out, g)
	}

	if err := rows.Err(); err != nil {
		logger.Error("rows iteration error", zap.Error(err))
		return nil, entity.ErrInternal
	}

	logger.Info("fetched game metadata",
		zap.Int("requested_ids", len(ids)),
		zap.Int("found_records", len(out)),
	)

	return out, nil
}
