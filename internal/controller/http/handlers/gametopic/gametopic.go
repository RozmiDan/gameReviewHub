package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GET   /games/{game_id}

type TopicGameGetter interface {
	GetTopicGame(ctx context.Context, gameID string) (*entity.Game, error)
}

func NewGameTopicHandler(baseLogger *zap.Logger, uc TopicGameGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1) забираем request_id из middleware и кладем в ctx
		reqID := middleware.GetReqID(r.Context())
		ctx := context.WithValue(r.Context(), entity.RequestIDKey{}, reqID)
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		// 2) оборачиваем логгер
		logger := baseLogger.With(zap.String("handler", "GameTopicHandler"), zap.String("request_id", reqID))

		gameID := chi.URLParam(r, "game_id")
		if _, err := uuid.Parse(gameID); err != nil {
			logger.Warn("invalid game_id", zap.String("game_id", gameID), zap.Error(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"invalid_game_id", "game_id is not a valid UUID"},
			})
			return
		}

		game, err := uc.GetTopicGame(ctx, gameID)
		if err != nil {
			// timeout
			if ctx.Err() == context.DeadlineExceeded {
				logger.Error("timeout exceeded", zap.Error(err))
				render.Status(r, http.StatusGatewayTimeout)
				render.JSON(w, r, ErrorResponse{
					Error: APIError{"timeout_exceeded", "request took longer than 2s"},
				})
				return
			}
			// игра не найдена
			if errors.Is(err, entity.ErrGameNotFound) {
				logger.Info("game not found", zap.String("game_id", gameID))
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, ErrorResponse{
					Error: APIError{"not_found", "game not found"},
				})
				return
			}
			// другая ошибка
			logger.Error("GetTopicGame failed", zap.Error(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"internal_error", "could not get game"},
			})
			return
		}

		// 5) форматируем ответ
		render.Status(r, http.StatusOK)
		render.JSON(w, r, GameTopicResponse{Data: *game})
	}
}
