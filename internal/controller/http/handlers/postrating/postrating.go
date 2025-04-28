package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

// POST  /games/{game_id}/rating

type RatingPoster interface {
	PostRating(ctx context.Context, gameID, userID string, rating int32) (bool, error)
}

type Response struct {
	Error  string `json:"error,omitempty"`
	Status string `json:"status"`
}

func NewRatingPostHandler(logger *zap.Logger, uc RatingPoster) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With(zap.String("func", "RatingPostHandler"),
			zap.String("request_id", middleware.GetReqID(r.Context())),
		)

	}
}
