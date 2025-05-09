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

// POST  /games/{game_id}/rating

type RatingPoster interface {
	PostRating(ctx context.Context, gameID, userID string, rating int32) error
}

// NewRatingPostHandler обрабатывает POST /games/{game_id}/rating — публикацию оценки игры.
// @Summary     Поставить оценку игре
// @Description Отправить новую оценку (1–10) для указанной игры
// @Tags        games
// @Accept      json
// @Produce     json
// @Param       game_id  path      string              true  "Идентификатор игры"
// @Param       payload  body      PostRatingRequest   true  "Тело запроса с user_id и rating"
// @Success     200      {object}  interface{}            "Пустой ответ — OK"
// @Failure     400      {object}  ErrorResponse       "Некорректный запрос"
// @Failure     404      {object}  ErrorResponse       "Игра не найдена"
// @Failure     503      {object}  ErrorResponse       "Брокер недоступен"
// @Failure     504      {object}  ErrorResponse       "Таймаут обработки запроса"
// @Failure     500      {object}  ErrorResponse       "Внутренняя ошибка сервера"
// @Router      /games/{game_id}/rating [post]
func NewRatingPostHandler(baseLogger *zap.Logger, uc RatingPoster) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1) Получаем request_id и создаём новый контекст с таймаутом
		reqID := chi.URLParam(r, "request_id") // если вы сохраняете его в URL-параметре
		if reqID == "" {
			reqID = middleware.GetReqID(r.Context())
		}
		ctx := context.WithValue(r.Context(), entity.RequestIDKey{}, reqID)
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		// 2) Оборачиваем логгер
		logger := baseLogger.
			With(zap.String("handler", "NewRatingPostHandler"), zap.String("request_id", reqID))

		// 3) Валидация game_id из URL
		gameID := chi.URLParam(r, "game_id")
		if _, err := uuid.Parse(gameID); err != nil {
			logger.Warn("invalid game_id", zap.String("game_id", gameID), zap.Error(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"invalid_game_id", "game_id is not a valid UUID"},
			})
			return
		}

		// 4) Декодируем тело
		var payload PostRatingRequest
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
				Error: APIError{"invalid_json", "cannot parse request body"},
			})
			return
		}

		// 5) Доп. валидация user_id и rating
		if _, err := uuid.Parse(payload.UserID); err != nil {
			logger.Warn("invalid user_id", zap.String("user_id", payload.UserID), zap.Error(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"invalid_user_id", "user_id is not a valid UUID"},
			})
			return
		}
		if payload.Rating < 1 || payload.Rating > 10 {
			logger.Warn("rating out of range", zap.Int32("rating", payload.Rating))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"invalid_rating", "rating must be between 1 and 10"},
			})
			return
		}

		// 6) Основная бизнес-логика
		if err := uc.PostRating(ctx, gameID, payload.UserID, payload.Rating); err != nil {
			switch {
			case errors.Is(err, entity.ErrBrokerUnavailable):
				logger.Error("broker unavailable", zap.Error(err))
				render.Status(r, http.StatusServiceUnavailable)
				render.JSON(w, r, ErrorResponse{
					Error: APIError{"service_unavailable", "unable to publish rating"},
				})
				return

			case errors.Is(err, entity.ErrGameNotFound):
				logger.Info("game not found", zap.String("game_id", gameID))
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, ErrorResponse{
					Error: APIError{"not_found", "game not found"},
				})
				return

			case errors.Is(err, context.DeadlineExceeded):
				logger.Error("timeout exceeded", zap.Error(err))
				render.Status(r, http.StatusGatewayTimeout)
				render.JSON(w, r, ErrorResponse{
					Error: APIError{"timeout_exceeded", "request took longer than 2 seconds"},
				})
				return
			}

			logger.Error("failed to post rating", zap.Error(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"internal error", "could not submit rating"},
			})
			return
		}

		// 7) Успех — пустой ответ 200 OK
		render.Status(r, http.StatusOK)
		render.JSON(w, r, struct{}{})
	}
}
