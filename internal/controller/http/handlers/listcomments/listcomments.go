package handlers

import (
	"context"
	"net/http"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

// GET  /games/{game_id}/comments?limit=&offset=

type ListCommentsGetter interface {
	GetListComments(ctx context.Context, gameID string, limit, offset int) ([]entity.Comment, error)
}

type Response struct {
	Error  string `json:"error,omitempty"`
	Status string `json:"status"`
}

func NewListCommentsHandler(logger *zap.Logger, uc ListCommentsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With(zap.String("func", "ListCommentsHandler"),
			zap.String("request_id", middleware.GetReqID(r.Context())),
		)

	}
}
