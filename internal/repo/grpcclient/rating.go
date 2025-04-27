package rating

import (
	"context"
	"time"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	ratingv1 "github.com/RozmiDan/gamehub-protos/gen/go/gamehub"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	api    ratingv1.RatingServiceClient
	logger *zap.Logger
}

func New(ctx context.Context, log *zap.Logger, addr string,
	timeout time.Duration) (*Client, error) {

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpc_zap.UnaryClientInterceptor(log),
			//grpc_prometheus.UnaryClientInterceptor,
		),
	)

	if err != nil {
		return nil, err
	}

	dialCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	conn.Connect()
	if !conn.WaitForStateChange(dialCtx, connectivity.Idle) {
		conn.Close()
		return nil, status.Error(codes.Unavailable, "could not connect to rating service")
	}

	api := ratingv1.NewRatingServiceClient(conn)
	return &Client{api: api, logger: log}, nil
}

func (c *Client) SubmitRating(ctx context.Context, userID, gameID string, rating int32) (bool, error) {
	resp, err := c.api.SubmitRating(ctx, &ratingv1.SubmitRatingRequest{
		UserId: userID,
		GameId: gameID,
		Rating: rating,
	})

	if err != nil {
		return false, err
	}

	if resp.Success {
		return true, nil
	} else {
		return false, nil
	}
}

func (c *Client) GetGameRating(ctx context.Context, gameID string) (*entity.GameRating, error) {
	resp, err := c.api.GetGameRating(ctx, &ratingv1.GetGameRatingRequest{
		GameId: gameID,
	})

	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, entity.ErrGameNotFound
		}
		return &entity.GameRating{}, err
	}

	respGame := &entity.GameRating{
		GameID:        resp.GetGameId(),
		RatingsCount:  resp.GetRatingsCount(),
		AverageRating: resp.GetAverageRating(),
	}

	return respGame, nil
}

func (c *Client) GetTopGames(ctx context.Context, limit, offset int32) ([]entity.GameRating, error) {
	resp, err := c.api.GetTopGames(ctx, &ratingv1.GetTopGamesRequest{
		Limit:  limit,
		Offset: offset,
	})

	topGames := []entity.GameRating{}

	// проверка ошибок на NotFound и т.д.

	if err != nil {
		return topGames, err
	}

	for _, it := range resp.Games {
		topGames = append(topGames, entity.GameRating{
			GameID:        it.GetGameId(),
			RatingsCount:  it.GetRatingsCount(),
			AverageRating: it.GetAverageRating(),
		})
	}

	return topGames, nil
}
