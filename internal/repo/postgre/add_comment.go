package postgres_storage

import (
	"context"
	"errors"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (r *RatingRepository) AddComment(ctx context.Context, gameID, userID, text string) (string, error) {
	// 1) забираем request_id
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)

	// 2) оборачиваем логгер
	logger := r.logger.With(zap.String("func", "AddComment"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	// 2) готовим и выполняем запрос
	const sqlQuery = `
        INSERT INTO comments(game_id, user_id, text)
        VALUES($1, $2, $3)
		RETURNING id
    `

	var commentID string
	err := r.pg.Pool.QueryRow(ctx, sqlQuery, gameID, userID, text).Scan(&commentID)

	if err != nil {
		// если ключ game_id не существует → 23503 foreign_key_violation
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			logger.Info("game_id not found", zap.String("game_id", gameID))
			return "", entity.ErrGameNotFound
		}
		logger.Error("failed to insert comment", zap.Error(err))
		return "", entity.ErrInsertComment
	}

	logger.Info("successfuly insert comment", zap.String("commentID", commentID))

	return commentID, nil
}
