package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	jsondecoder "github.com/RozmiDan/gameReviewHub/pkg/json_decoder"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// POST /games

type GameTopicCreator interface {
	CreateGameTopic(ctx context.Context, game *entity.Game) (string, error)
}

func NewCreateGameHandler(baseLogger *zap.Logger, uc GameTopicCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1) Получаем request_id и создаём новый контекст с таймаутом
		reqID := chi.URLParam(r, "request_id")
		if reqID == "" {
			reqID = middleware.GetReqID(r.Context())
		}
		ctx := context.WithValue(r.Context(), entity.RequestIDKey{}, reqID)
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		// 2) Оборачиваем логгер
		logger := baseLogger.With(zap.String("handler", "CreateGameHandler"), zap.String("request_id", reqID))

		// 3) Декодируем тело
		var payload CreateGameRequest
		if err := jsondecoder.DecodeJSONBody(w, r, &payload); err != nil {
			mr, ok := err.(*jsondecoder.MalformedRequest)
			if ok {
				logger.Warn("malformed request body", zap.Error(err))
				render.Status(r, mr.Status)
				render.JSON(w, r, ErrorResponse{
					Error: APIError{mr.Msg, mr.Msg},
				})
				return
			}
			logger.Error("failed to decode JSON", zap.Error(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"missing_fields", "cannot parse request body"},
			})
			return
		}

		if payload.ID != "" {
			if _, err := uuid.Parse(payload.ID); err != nil {
				logger.Warn("invalid game_id", zap.String("game_id", payload.ID), zap.Error(err))
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, ErrorResponse{
					Error: APIError{"invalid_game_id", "game_id is not a valid UUID"},
				})
				return
			}
		}

		if payload.Creator == "" || payload.Description == "" ||
			payload.Genre == "" || payload.Name == "" {
			logger.Warn("missing required fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"invalid_text", ""},
			})
			return
		}

		releaseDate, err := time.Parse("2006-01-02", payload.ReleaseDate)
		if err != nil {
			logger.Warn("invalid release_date format", zap.String("value", payload.ReleaseDate))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"invalid_date", "release_date must be YYYY-MM-DD"},
			})
			return
		}

		// 6) Основная бизнес-логика
		commentID, err := uc.CreateGameTopic(ctx, &entity.Game{
			ID:          payload.ID,
			Name:        payload.Name,
			Genre:       payload.Genre,
			Creator:     payload.Creator,
			Description: payload.Description,
			ReleaseDate: releaseDate,
		})

		switch {
		case errors.Is(err, entity.ErrInsertGame):
			logger.Error("failed to insert game", zap.Error(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"insert_failed", "could not create game"},
			})
			return

		case errors.Is(err, entity.ErrGameAlreadyExists):
			logger.Info("duplicate game", zap.String("name", payload.Name))
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"already_exists", "game with this name already exists"},
			})
			return

		case errors.Is(err, entity.ErrInternal):
			logger.Error("unexpected error adding game", zap.Error(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"internal_error", "internal server error"},
			})
			return

		case ctx.Err() == context.DeadlineExceeded:
			logger.Error("timeout creating game topic", zap.Error(err))
			render.Status(r, http.StatusGatewayTimeout)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"timeout_exceeded", "request took longer than 2 seconds"},
			})
			return

		case err != nil:
			logger.Error("unknown error", zap.Error(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"unknown_error", "could not create game"},
			})
			return
		}

		// 7) успех — 201 Created + Location + тело CreateGameResponse
		w.Header().Set("Location", "/games/"+commentID)
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, CreateGameResponse{ID: commentID})
	}
}
