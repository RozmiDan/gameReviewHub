// internal/usecase/post_rating_test.go
package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type fakeRepo struct {
	result []entity.GameInList
	err    error
}

func (f *fakeRepo) GetGameInfo(ctx context.Context, ids []string) ([]entity.GameInList, error) {
	return f.result, f.err
}

func (f *fakeRepo) GetGameTopic(context.Context, string) (*entity.Game, error) {
	panic("not used")
}
func (f *fakeRepo) GetCommentsGame(context.Context, string, int32, int32) ([]entity.Comment, error) {
	panic("not used")
}
func (f *fakeRepo) AddComment(context.Context, string, string, string) (string, error) {
	panic("not used")
}
func (f *fakeRepo) AddGameTopic(context.Context, *entity.Game) (string, error) {
	panic("not used")
}
func (f *fakeRepo) GetGameInfoBatch(context.Context, []string) ([]entity.GameInList, error) {
	return f.GetGameInfo(context.Background(), nil)
}

type fakeProducer struct {
	called bool
	msg    entity.RatingMessage
	err    error
}

func (f *fakeProducer) PublishRating(ctx context.Context, m entity.RatingMessage) error {
	f.called = true
	f.msg = m
	return f.err
}

func TestUsecase_PostRating(t *testing.T) {
	const (
		gameID = "game-123"
		userID = "user-xyz"
		rate   = 5
	)

	tests := []struct {
		name        string
		repoResult  []entity.GameInList
		repoErr     error
		producerErr error
		wantErr     error
		wantCalled  bool
	}{
		{
			name:       "repository internal error",
			repoErr:    entity.ErrInternal,
			wantErr:    entity.ErrInternal,
			wantCalled: false,
		},
		{
			name:       "repository unexpected error",
			repoErr:    errors.New("db down"),
			wantErr:    entity.ErrUnidentified,
			wantCalled: false,
		},
		{
			name:       "game not found",
			repoResult: []entity.GameInList{},
			wantErr:    entity.ErrGameNotFound,
			wantCalled: false,
		},
		{
			name:        "producer unavailable",
			repoResult:  []entity.GameInList{{ID: gameID}},
			producerErr: errors.New("kafka down"),
			wantErr:     entity.ErrBrokerUnavailable,
			wantCalled:  true,
		},
		{
			name:       "happy path",
			repoResult: []entity.GameInList{{ID: gameID}},
			wantErr:    nil,
			wantCalled: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			logger := zap.NewNop()

			repo := &fakeRepo{
				result: tc.repoResult,
				err:    tc.repoErr,
			}
			prod := &fakeProducer{err: tc.producerErr}

			uc := New(
				nil,
				repo,
				logger,
				prod,
				nopCache,
			)

			err := uc.PostRating(ctx, gameID, userID, rate)
			if tc.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.wantErr)
			}

			require.Equal(t, tc.wantCalled, prod.called, "PublishRating called")
			if tc.wantCalled && tc.wantErr == nil {
				require.Equal(t, entity.RatingMessage{
					GameID: gameID,
					UserID: userID,
					Rating: rate,
				}, prod.msg)
			}
		})
	}
}
