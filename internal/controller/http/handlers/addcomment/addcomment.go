package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

// POST /games/{game_id}/comments

type CommentPoster interface {
	AddComment(ctx context.Context, gameID, userID, text string) error
}

type Response struct {
	Error  string `json:"error,omitempty"`
	Status string `json:"status"`
}

func NewAddCommentHandler(logger *zap.Logger, uc CommentPoster) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With(zap.String("func", "AddCommentHandler"),
			zap.String("request_id", middleware.GetReqID(r.Context())),
		)

	}
}
