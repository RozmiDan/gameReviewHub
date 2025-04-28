package postgres_storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	postgres_storage "github.com/RozmiDan/gameReviewHub/pkg/postgres"
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

func (r *RatingRepository) GetGameTopic(ctx context.Context, ids string) (*entity.Game, error) {
	return &entity.Game{}, nil
}

func (r *RatingRepository) GetGameInfo(ctx context.Context, ids []string) ([]entity.GameInList, error) {
	logger := r.logger.With(zap.String("func", "GetGameInfo"), zap.Int("count", len(ids)))
	if len(ids) == 0 {
		return nil, nil
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

	logger.Info("games successfuly found", zap.Int("count", len(out)))

	return out, nil
}

func (r *RatingRepository) GetCommentsGame(ctx context.Context, gameID string, limit, offset int) ([]entity.Comment, error) {
	return []entity.Comment{}, nil
}

func (r *RatingRepository) AddComment(ctx context.Context, gameID, userID, text string) (string, error) {
	return "", nil
}
