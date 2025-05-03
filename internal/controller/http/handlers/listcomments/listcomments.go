package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GET  /games/{game_id}/comments?limit=&offset=

type ListCommentsGetter interface {
	GetListComments(ctx context.Context, gameID string, limit, offset int32) ([]entity.Comment, error)
}

func NewListCommentsHandler(baseLogger *zap.Logger, uc ListCommentsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1) request_id и таймаут
		reqID := middleware.GetReqID(r.Context())
		ctx := context.WithValue(r.Context(), entity.RequestIDKey{}, reqID)
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		// 2) оборачиваем логгер
		logger := baseLogger.With(zap.String("handler", "ListCommentsHandler"), zap.String("request_id", reqID))

		// 3) валидируем game_id из URL
		gameID := chi.URLParam(r, "game_id")
		if _, err := uuid.Parse(gameID); err != nil {
			logger.Warn("invalid game_id", zap.String("game_id", gameID), zap.Error(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"invalid_game_id", "game_id must be a valid UUID"},
			})
			return
		}

		// 4) парсим и валидируем limit/offset
		limit, offset, errStruct := parseAndValidatePaging(r, logger)
		if errStruct != nil {
			render.JSON(w, r, errStruct)
			return
		}

		// 5) вызываем бизнес-логику
		comments, err := uc.GetListComments(ctx, gameID, limit, offset)
		if err != nil {
			switch {
			case errors.Is(err, entity.ErrTimeout):
				logger.Error("timeout fetching comments", zap.Error(err))
				render.Status(r, http.StatusGatewayTimeout)
				render.JSON(w, r, ErrorResponse{
					Error: APIError{"timeout_exceeded", "request took longer than 2s"},
				})
			case errors.Is(err, entity.ErrInternal):
				logger.Error("internal error fetching comments", zap.Error(err))
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, ErrorResponse{
					Error: APIError{"internal_error", "could not fetch comments"},
				})
			default:
				logger.Error("unexpected error fetching comments", zap.Error(err))
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, ErrorResponse{
					Error: APIError{"internal_error", "could not fetch comments"},
				})
			}
			return
		}
		// 6) формируем и отдаем ответ
		resp := ListCommentsResponse{
			Data: comments,
			Meta: &Pagination{
				Limit:  limit,
				Offset: offset,
				Count:  len(comments),
			},
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}

// читаем limit/offset, логируем и возвращаем ErrorResponse, если невалидно
func parseAndValidatePaging(r *http.Request, logger *zap.Logger) (limit, offset int32, errResp *ErrorResponse) {
	q := r.URL.Query()

	limit = 10
	if s := q.Get("limit"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			limit = int32(v)
		} else {
			logger.Warn("invalid limit param", zap.String("limit", s), zap.Error(err))
			render.Status(r, http.StatusBadRequest)
			errResp = &ErrorResponse{Error: APIError{"invalid_limit", "limit must be a positive integer"}}
			return
		}
	}

	offset = 0
	if s := q.Get("offset"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			offset = int32(v)
		} else {
			logger.Warn("invalid offset param", zap.String("offset", s), zap.Error(err))
			render.Status(r, http.StatusBadRequest)
			errResp = &ErrorResponse{Error: APIError{"invalid_offset", "offset must be a non-negative integer"}}
			return
		}
	}

	if limit <= 0 {
		logger.Warn("limit out of range", zap.Int32("limit", limit))
		render.Status(r, http.StatusBadRequest)
		errResp = &ErrorResponse{Error: APIError{"invalid_limit", "limit must be > 0"}}
		return
	}
	if offset < 0 {
		logger.Warn("offset out of range", zap.Int32("offset", offset))
		render.Status(r, http.StatusBadRequest)
		errResp = &ErrorResponse{Error: APIError{"invalid_offset", "offset must be >= 0"}}
		return
	}

	return
}
