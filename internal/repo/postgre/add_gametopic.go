package postgres_storage

import (
	"context"
	"errors"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (r *RatingRepository) AddGameTopic(ctx context.Context, gameInfo *entity.Game) (string, error) {
	// 1) забираем request_id
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)

	// 2) оборачиваем логгер
	logger := r.logger.With(zap.String("func", "AddGameTopic"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	var args []interface{}
	var sqlQuery string

	if gameInfo.ID != "" {
		sqlQuery = `
        INSERT INTO games (id, name, genre, creator, description, release_date)
          VALUES ($1, $2, $3, $4, $5, $6)
          RETURNING id;
    	`

		args = []interface{}{
			gameInfo.ID,
			gameInfo.Name,
			gameInfo.Genre,
			gameInfo.Creator,
			gameInfo.Description,
			gameInfo.ReleaseDate,
		}

	} else {
		// 2) готовим и выполняем запрос
		sqlQuery = `
			INSERT INTO games (name, genre, creator, description, release_date)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (name) DO NOTHING
			RETURNING id;
		`

		args = []interface{}{
			gameInfo.Name,
			gameInfo.Genre,
			gameInfo.Creator,
			gameInfo.Description,
			gameInfo.ReleaseDate,
		}
	}

	var gameID string
	err := r.pg.Pool.QueryRow(ctx, sqlQuery, args...).Scan(&gameID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Info("game already exists, skipping insert",
				zap.String("name", gameInfo.Name),
			)
			return "", entity.ErrGameAlreadyExists
		}
		// ловим любое другое pg-ошибку
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// unique_violation
			switch pgErr.Code {
			case "23505": // unique_violation
				return "", entity.ErrGameAlreadyExists
			case "23503": // foreign_key_violation (unlikely here)
				logger.Warn("foreign key violation", zap.Error(err))
				return "", entity.ErrInternal
			}
		}
		logger.Error("failed to insert game", zap.Error(err))
		return "", entity.ErrInsertGame
	}

	logger.Info("successfuly insert game", zap.String("gameID", gameID))

	return gameID, nil
}
