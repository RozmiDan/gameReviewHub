package handlers

import (
	"context"
	"net/http"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

// GET   /games/{game_id}

type TopicGameGetter interface {
	GetTopicGame(ctx context.Context, gameID string) (entity.Game, error)
}

type Response struct {
	Error  string `json:"error,omitempty"`
	Status string `json:"status"`
}

func NewGameTopicHandler(logger *zap.Logger, uc TopicGameGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With(zap.String("func", "GameTopicHandler"),
			zap.String("request_id", middleware.GetReqID(r.Context())),
		)

	}
}
