package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockRepo struct {
	returnID string
	returnErr error
}

func (m *mockRepo) GetGameTopic(ctx context.Context, gameID string) (*entity.Game, error) {
	panic("not implemented")
}
func (m *mockRepo) GetGameInfo(ctx context.Context, ids []string) ([]entity.GameInList, error) {
	panic("not implemented")
}
func (m *mockRepo) GetCommentsGame(ctx context.Context, gameID string, limit, offset int32) ([]entity.Comment, error) {
	panic("not implemented")
}
func (m *mockRepo) AddComment(ctx context.Context, gameID, userID, text string) (string, error) {
	return m.returnID, m.returnErr
}
func (m *mockRepo) AddGameTopic(ctx context.Context, gameInfo *entity.Game) (string, error) {
	panic("not implemented")
}

var nopRatingClient RatingClient = nil
var nopProducer RatingProducer = nil

func TestAddComment(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()

	const (
		gameID = "game-123"
		userID = "user-456"
		text   = "hello"
	)

	tests := []struct {
		name        string
		repoErr     error
		repoID      string
		wantID      string
		wantErr     error
	}{
		{
			name:    "game not found",
			repoErr: entity.ErrGameNotFound,
			wantErr: entity.ErrGameNotFound,
		},
		{
			name:    "insert comment failure",
			repoErr: entity.ErrInsertComment,
			wantErr: entity.ErrInsertComment,
		},
		{
			name:    "unexpected repo error",
			repoErr: errors.New("db is down"),
			wantErr: entity.ErrInternal,
		},
		{
			name:    "success",
			repoErr: nil,
			repoID:  "comment-789",
			wantID:  "comment-789",
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uc := New(
				nopRatingClient,
				&mockRepo{returnID: tc.repoID, returnErr: tc.repoErr},
				logger,
				nopProducer,
			)

			gotID, gotErr := uc.AddComment(ctx, gameID, userID, text)
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