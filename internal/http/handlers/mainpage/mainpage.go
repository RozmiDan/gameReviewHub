package handlers

import (
	"context"
	"net/http"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type GamesListGetter interface {
	GetListGames(ctx context.Context, limit, offset int) ([]entity.GameRating, error)
}

type Response struct {
	Error  string `json:"error,omitempty"`
	Status string `json:"status"`
}

func NewMainpageHandler(logger *zap.Logger, uc GamesListGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With(zap.String("func", "MainpageHandler"),
			zap.String("request_id", middleware.GetReqID(r.Context())),
		)

	}
}
