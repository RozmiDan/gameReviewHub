package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type fakeRatingClient struct {
	topGames []entity.GameRating
	err      error
}

func (f *fakeRatingClient) SubmitRating(ctx context.Context, userID, gameID string, rating int32) (bool, error) {
	return false, nil
}
func (f *fakeRatingClient) GetGameRating(ctx context.Context, gameID string) (*entity.GameRating, error) {
	return nil, nil
}
func (f *fakeRatingClient) GetTopGames(ctx context.Context, limit, offset int32) ([]entity.GameRating, error) {
	return f.topGames, f.err
}

type fakeGameRepo struct {
	metas []entity.GameInList
	err   error
}

func (f *fakeGameRepo) GetGameTopic(ctx context.Context, gameID string) (*entity.Game, error) {
	return nil, nil
}
func (f *fakeGameRepo) GetGameInfo(ctx context.Context, ids []string) ([]entity.GameInList, error) {
	return f.metas, f.err
}
func (f *fakeGameRepo) GetCommentsGame(ctx context.Context, gameID string, limit, offset int32) ([]entity.Comment, error) {
	return nil, nil
}
func (f *fakeGameRepo) AddComment(ctx context.Context, gameID, userID, text string) (string, error) {
	return "", nil
}
func (f *fakeGameRepo) AddGameTopic(ctx context.Context, gameInfo *entity.Game) (string, error) {
	return "", nil
}

func TestGetListGames(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()

	tests := []struct {
		name        string
		ratings     []entity.GameRating
		ratingErr   error
		metas       []entity.GameInList
		metaErr     error
		expectedOut []entity.GameInList
		expectedErr bool
	}{
		{
			name:        "rating service error",
			ratings:     nil,
			ratingErr:   errors.New("rpc failed"),
			expectedErr: true,
		},
		{
			name:        "empty ratings",
			ratings:     []entity.GameRating{},
			ratingErr:   nil,
			expectedOut: []entity.GameInList{},
		},
		{
			name:        "db metadata error",
			ratings:     []entity.GameRating{{GameID: "g1", AverageRating: 5.5}},
			ratingErr:   nil,
			metas:       nil,
			metaErr:     errors.New("db down"),
			expectedErr: true,
		},
		{
			name:    "missing metadata for some games",
			ratings: []entity.GameRating{{GameID: "g1", AverageRating: 5.5}, {GameID: "g2", AverageRating: 3.0}},
			metas:   []entity.GameInList{{ID: "g1", Name: "One", Genre: "A"}},
			expectedOut: []entity.GameInList{
				{ID: "g1", Name: "One", Genre: "A", Rating: 5.5},
			},
		},
		{
			name:    "happy path",
			ratings: []entity.GameRating{{GameID: "g1", AverageRating: 5.5}, {GameID: "g2", AverageRating: 3.0}},
			metas: []entity.GameInList{
				{ID: "g2", Name: "Two", Genre: "B"},
				{ID: "g1", Name: "One", Genre: "A"},
			},
			expectedOut: []entity.GameInList{
				{ID: "g1", Name: "One", Genre: "A", Rating: 5.5},
				{ID: "g2", Name: "Two", Genre: "B", Rating: 3.0},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uc := New(
				&fakeRatingClient{topGames: tc.ratings, err: tc.ratingErr},
				&fakeGameRepo{metas: tc.metas, err: tc.metaErr},
				logger,
				nil,
			)

			out, err := uc.GetListGames(ctx, 10, 0)
			if tc.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}
