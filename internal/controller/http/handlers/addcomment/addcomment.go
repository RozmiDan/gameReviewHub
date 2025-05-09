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

// POST /games/{game_id}/comments

type CommentPoster interface {
	AddComment(ctx context.Context, gameID, userID, text string) (string, error)
}

// AddCommentHandler добавляет новый комментарий к игре.
// @Summary     Постинг комментария
// @Description Добавляет комментарий пользователя к указанной игре.
// @Tags        comments
// @Accept      json
// @Produce     json
// @Param       game_id   path     string             true  "UUID игры"
// @Param       body      body     PostCommentRequest true  "Тело запроса с полем user_id и text"
// @Success     200       {object} AddCommentResponse  "ID созданного комментария"
// @Failure     400       {object} ErrorResponse        "Некорректные входные данные"
// @Failure     404       {object} ErrorResponse        "Игра не найдена"
// @Failure     504       {object} ErrorResponse        "Таймаут запроса"
// @Failure     500       {object} ErrorResponse        "Внутренняя ошибка сервера"
// @Router      /games/{game_id}/comments [post]
func NewAddCommentHandler(baseLogger *zap.Logger, uc CommentPoster) http.HandlerFunc {
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
		logger := baseLogger.With(zap.String("handler", "AddCommentHandler"), zap.String("request_id", reqID))

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
		var payload PostCommentRequest
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
		if len(payload.Text) == 0 || len(payload.Text) > 1000 {
			logger.Warn("invalid text length", zap.Int("text_size", len(payload.Text)))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"invalid_text", "comment size must be between 0 and 1000"},
			})
			return
		}

		// 6) Основная бизнес-логика
		commentID, err := uc.AddComment(ctx, gameID, payload.UserID, payload.Text)
		switch {
		case errors.Is(err, entity.ErrGameNotFound):
			logger.Info("game not found", zap.String("game_id", gameID))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"not_found", "game not found"},
			})
			return

		case errors.Is(err, entity.ErrInsertComment):
			logger.Error("failed to insert comment", zap.String("game_id", gameID), zap.Error(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"insert_failed", "could not create comment"},
			})
			return

		case errors.Is(err, entity.ErrInternal):
			logger.Error("unexpected error adding comment", zap.Error(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"internal_error", "internal server error"},
			})
			return

		case ctx.Err() == context.DeadlineExceeded:
			logger.Error("timeout adding comment", zap.Error(err))
			render.Status(r, http.StatusGatewayTimeout)
			render.JSON(w, r, ErrorResponse{
				Error: APIError{"timeout_exceeded", "request took longer than 2 seconds"},
			})
			return

		default:
			// если err == nil, продолжаем
		}

		// 7) Отдаем ID нового комментария
		render.Status(r, http.StatusOK)
		render.JSON(w, r, AddCommentResponse{
			ID: commentID,
		})
	}
}
