package postgres_storage

import (
	"context"
	"errors"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (r *RatingRepository) GetGameTopic(ctx context.Context, gameID string) (*entity.Game, error) {
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)
	logger := r.logger.With(zap.String("func", "GetGameTopic"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	// 2) готовим и выполняем запрос
	const sqlQuery = `
        SELECT name, genre, creator, description, release_date
        FROM games
        WHERE id = $1
    `
	row := r.pg.Pool.QueryRow(ctx, sqlQuery, gameID)

	// 3) сканируем результат
	g := &entity.Game{}
	g.ID = gameID
	if err := row.Scan(
		&g.Name,
		&g.Genre,
		&g.Creator,
		&g.Description,
		&g.ReleaseDate,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// если игра не найдена — возвращаем ошибку
			logger.Info("game not found", zap.String("game_id", gameID))
			return &entity.Game{}, entity.ErrGameNotFound
		}
		// любая другая ошибка — логируем и прокидываем
		logger.Error("failed to scan game record", zap.Error(err))
		return &entity.Game{}, err
	}

	// 4) всё успешно — возвращаем объект
	logger.Info("game fetched from database", zap.String("game_id", gameID))
	return g, nil
}
