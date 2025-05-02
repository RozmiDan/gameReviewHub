package postgres_storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	postgres_storage "github.com/RozmiDan/gameReviewHub/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type RatingRepository struct {
	pg     *postgres_storage.Postgres
	logger *zap.Logger
}

func New(pg *postgres_storage.Postgres, logger *zap.Logger) *RatingRepository {
	logger = logger.With(zap.String("layer", "MainRepository"))
	return &RatingRepository{pg, logger}
}

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

func (r *RatingRepository) GetGameInfo(ctx context.Context, ids []string) ([]entity.GameInList, error) {

	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)
	logger := r.logger
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}
	logger = logger.With(zap.String("func", "GetGameInfo"))

	if len(ids) == 0 {
		logger.Debug("no ids provided, returning empty list")
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
		return nil, err
	}
	defer rows.Close()

	var out []entity.GameInList
	for rows.Next() {
		var g entity.GameInList
		if err := rows.Scan(&g.ID, &g.Name, &g.Genre); err != nil {
			logger.Error("scan failed", zap.Error(err))
			return nil, err
		}
		out = append(out, g)
	}

	logger.Info("fetched game metadata",
		zap.Int("requested_ids", len(ids)),
		zap.Int("found_records", len(out)),
	)

	return out, nil
}

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
	rows, err := r.pg.Pool.Query(ctx, sqlQuery, gameID, limit, offset)
	if err != nil {
		logger.Error("query failed", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	// 3) сканируем результат
	var comments []entity.Comment
	for rows.Next() {
		var comment entity.Comment
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.Text, &comment.CreatedAt); err != nil {
			logger.Error("scan failed", zap.Error(err))
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
        logger.Error("rows iteration error", zap.Error(err))
        return nil, err
    }

	logger.Info("fetched comments", zap.Int("found_records", len(comments)))

	return comments, nil
}

func (r *RatingRepository) AddComment(ctx context.Context, gameID, userID, text string) (string, error) {
	return "", nil
}
