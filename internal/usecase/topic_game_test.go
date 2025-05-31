package usecase

import (
	"context"
	"testing"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type fakeTopicRepo struct {
	game *entity.Game
	err  error
}

func (f *fakeTopicRepo) GetGameTopic(ctx context.Context, id string) (*entity.Game, error) {
	return f.game, f.err
}

func (f *fakeTopicRepo) GetGameInfo(ctx context.Context, ids []string) ([]entity.GameInList, error) {
	panic("not used")
}
func (f *fakeTopicRepo) GetCommentsGame(ctx context.Context, gameID string, limit, offset int32) ([]entity.Comment, error) {
	panic("not used")
}
func (f *fakeTopicRepo) AddComment(ctx context.Context, gameID, userID, text string) (string, error) {
	panic("not used")
}
func (f *fakeTopicRepo) AddGameTopic(ctx context.Context, game *entity.Game) (string, error) {
	panic("not used")
}

type fakeratingClient struct {
	rating *entity.GameRating
	err    error
}

func (f *fakeratingClient) GetGameRating(ctx context.Context, gameID string) (*entity.GameRating, error) {
	return f.rating, f.err
}

func (f *fakeratingClient) SubmitRating(ctx context.Context, userID, gameID string, rating int32) (bool, error) {
	panic("not used")
}
func (f *fakeratingClient) GetTopGames(ctx context.Context, limit, offset int32) ([]entity.GameRating, error) {
	panic("not used")
}

func TestUsecase_GetTopicGame(t *testing.T) {
	const gid = "game-1"

	emptyGame := &entity.Game{}

	tests := []struct {
		name      string
		repoGame  *entity.Game
		repoErr   error
		rating    *entity.GameRating
		ratingErr error
		wantGame  *entity.Game
		wantErr   error
	}{
		{
			name:     "repo: not found",
			repoErr:  entity.ErrGameNotFound,
			wantGame: emptyGame, wantErr: entity.ErrGameNotFound,
		},
		{
			name:      "no ratings yet",
			repoGame:  &entity.Game{ID: gid, Name: "G"},
			ratingErr: entity.ErrGameNotFound,
			wantGame:  &entity.Game{ID: gid, Name: "G"}, wantErr: nil,
		},
		{
			name:      "invalid uuid from rating client",
			repoGame:  &entity.Game{ID: gid},
			ratingErr: entity.ErrInvalidUUID,
			wantGame:  nil, wantErr: entity.ErrInvalidUUID,
		},
		{
			name:      "service unavailable",
			repoGame:  &entity.Game{ID: gid},
			ratingErr: entity.ErrServiceUnavailable,
			wantGame:  &entity.Game{ID: gid}, wantErr: nil,
		},
		{
			name:      "internal rating error",
			repoGame:  &entity.Game{ID: gid},
			ratingErr: entity.ErrInternalRating,
			wantGame:  &entity.Game{ID: gid}, wantErr: nil,
		},
		{
			name:     "happy path",
			repoGame: &entity.Game{ID: gid, Name: "G"},
			rating:   &entity.GameRating{GameID: gid, AverageRating: 4.5, RatingsCount: 10},
			wantGame: &entity.Game{ID: gid, Name: "G", Rating: entity.GameRating{GameID: gid, AverageRating: 4.5, RatingsCount: 10}},
			wantErr:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeTopicRepo{game: tc.repoGame, err: tc.repoErr}
			rcl := &fakeratingClient{rating: tc.rating, err: tc.ratingErr}
			uc := New(rcl, repo, zap.NewNop(), nil, nopCache)

			got, err := uc.GetTopicGame(context.Background(), gid)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.wantGame, got)
		})
	}
}
