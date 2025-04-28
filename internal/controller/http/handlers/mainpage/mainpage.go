package handlers

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

// 1) GET  /games?limit=&offset=

type GamesListGetter interface {
	GetListGames(ctx context.Context, limit, offset int32) ([]entity.GameInList, error)
}

type RequestIDKey struct{}

func NewMainpageHandler(logger *zap.Logger, uc GamesListGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())

		ctx := context.WithValue(r.Context(), RequestIDKey{}, reqID)
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		logger = logger.With(zap.String("func", "MainpageHandler"),
			zap.String("request_id", reqID),
		)

		q := r.URL.Query()

		limit, offset := parseLimitOffset(&q)

		list, err := uc.GetListGames(ctx, limit, offset)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				logger.Error("timeout exceeded", zap.Error(err))

				render.Status(r, http.StatusGatewayTimeout)
				render.JSON(w, r, ErrorResponse{
					Error: APIError{
						Code:    "timeout_exceeded",
						Message: "Request took longer than 2s",
					},
				})
			}

			logger.Error("cant get list games", zap.Error(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{
					Code:    "internal_error",
					Message: "Something went wrong",
				},
			})
			return
		}

		logger.Info("querry param", zap.Any("limit", limit), zap.Any("offset", offset))

		resp := ListGamesResponse{
			Data: list,
			Meta: &Pagination{
				Limit:  limit,
				Offset: offset,
				Count:  len(list),
			},
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}

func parseLimitOffset(q *url.Values) (limit, offset int32) {
	limit = 10
	if s := q.Get("limit"); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			limit = int32(n)
		}
	}
	offset = 0
	if s := q.Get("offset"); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			offset = int32(n)
		}
	}
	return
}
