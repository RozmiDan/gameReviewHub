package postgres_storage

import (
	"context"
	"errors"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (r *RatingRepository) GetCommentsGame(ctx context.Context, gameID string, limit, offset int32) ([]entity.Comment, error) {
	// 1) забираем request_id
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)

	// 2) оборачиваем логгер
	logger := r.logger.With(zap.String("func", "GetCommentsGame"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	// 2) готовим и выполняем запрос
	const sqlQuery = `
        SELECT id, user_id, text, created_at
        FROM comments
        WHERE game_id = $1
		ORDER BY created_at DESC
      	LIMIT $2 OFFSET $3
    `

	rows, err := r.pg.Pool.Query(ctx, sqlQuery, gameID, limit, offset*limit)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Invalid text representation, foreign key и тп
			logger.Error("postgres error querying comments", zap.Error(err))
			return nil, entity.ErrInternalComments
		}
		logger.Error("query failed", zap.Error(err))
		return nil, entity.ErrInternalComments
	}
	defer rows.Close()

	// 3) сканируем результат
	var comments []entity.Comment
	for rows.Next() {
		var comment entity.Comment
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.Text, &comment.CreatedAt); err != nil {
			logger.Error("scan failed", zap.Error(err))
			return nil, entity.ErrInternalComments
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		logger.Error("rows iteration error", zap.Error(err))
		return nil, entity.ErrInternalComments
	}

	logger.Info("fetched comments", zap.Int("found_records", len(comments)))

	return comments, nil
}
