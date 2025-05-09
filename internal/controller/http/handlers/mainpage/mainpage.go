package handlers

import (
	"context"
	"net/http"
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

// ListGamesHandler возвращает список игр с пагинацией.
// @Summary     Получить список игр
// @Description Возвращает упорядоченный по id список игр с поддержкой limit/offset.
// @Tags        games
// @Accept      json
// @Produce     json
// @Param       limit   query     int  false  "Максимальное число игр"       default(10)
// @Param       offset  query     int  false  "Смещение для пагинации"       default(0)
// @Success     200     {object}  ListGamesResponse   "Список игр и мета"
// @Failure     400     {object}  ErrorResponse       "Неверные параметры запроса"
// @Failure     504     {object}  ErrorResponse       "Таймаут обработки запроса"
// @Failure     500     {object}  ErrorResponse       "Внутренняя ошибка сервера"
// @Router      /games [get]
func NewMainpageHandler(baseLogger *zap.Logger, uc GamesListGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1) забираем request_id из middleware и кладем в ctx
		reqID := middleware.GetReqID(r.Context())
		ctx := context.WithValue(r.Context(), entity.RequestIDKey{}, reqID)
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		// 2) оборачиваем логгер
		logger := baseLogger.
			With(zap.String("handler", "MainpageHandler"), zap.String("request_id", reqID))

		// 3) парсим и валидируем query-параметры
		limit, offset, errStruct := parseAndValidatePaging(r, logger)
		if errStruct != nil {
			render.JSON(w, r, errStruct)
			return
		}

		// 4) вызываем usecase
		list, err := uc.GetListGames(ctx, limit, offset)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				logger.Error("timeout exceeded", zap.Error(err))
				render.Status(r, http.StatusGatewayTimeout)
				render.JSON(w, r, ErrorResponse{
					Error: APIError{"timeout_exceeded", "request took longer than 2s"},
				})
				return
			}
			logger.Error("GetListGames failed", zap.Error(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"internal_error", "could not fetch games list"},
			})
			return
		}

		// 5) форматируем ответ
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
