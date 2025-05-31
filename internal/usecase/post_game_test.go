// internal/usecase/create_game_topic_test.go
package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// mockRepo для AddGameTopic; остальные методы паникуют, если их вызовут
type mockGameRepo struct {
	returnID  string
	returnErr error
}

func (m *mockGameRepo) GetGameTopic(ctx context.Context, gameID string) (*entity.Game, error) {
	panic("not implemented")
}
func (m *mockGameRepo) GetGameInfo(ctx context.Context, ids []string) ([]entity.GameInList, error) {
	panic("not implemented")
}
func (m *mockGameRepo) GetCommentsGame(ctx context.Context, gameID string, limit, offset int32) ([]entity.Comment, error) {
	panic("not implemented")
}
func (m *mockGameRepo) AddComment(ctx context.Context, gameID, userID, text string) (string, error) {
	panic("not implemented")
}
func (m *mockGameRepo) AddGameTopic(ctx context.Context, game *entity.Game) (string, error) {
	return m.returnID, m.returnErr
}

func TestCreateGameTopic(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()

	// тестовая игра
	testGame := &entity.Game{
		Name:        "Test Game",
		Genre:       "GenreX",
		Creator:     "AuthorY",
		Description: "DescZ",
	}

	cases := []struct {
		name       string
		repoID     string
		repoErr    error
		wantID     string
		wantErr    error
	}{
		{
			name:    "already exists",
			repoErr: entity.ErrGameAlreadyExists,
			wantErr: entity.ErrGameAlreadyExists,
		},
		{
			name:    "insert failure",
			repoErr: entity.ErrInsertGame,
			wantErr: entity.ErrInsertGame,
		},
		{
			name:    "unexpected failure",
			repoErr: errors.New("db down"),
			wantErr: entity.ErrInternal,
		},
		{
			name:    "success",
			repoID:  "game-123",
			wantID:  "game-123",
			wantErr: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			uc := New(
				nopRatingClient,
				&mockGameRepo{returnID: tc.repoID, returnErr: tc.repoErr},
				logger,
				nopProducer,
				nopCache,
			)

			gotID, gotErr := uc.CreateGameTopic(ctx, testGame)
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, gotErr)
				assert.Empty(t, gotID)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, tc.wantID, gotID)
			}
		})
	}
}
