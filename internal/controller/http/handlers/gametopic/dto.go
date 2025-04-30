package handlers

import "github.com/RozmiDan/gameReviewHub/internal/entity"

// ListGamesResponse — обёртка для GET /games
type GameTopicResponse struct {
	Data entity.Game `json:"data"`
}

// --------------- ответы с ошибкой ---------------

// APIError — структура описания ошибки
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse — обёртка для не-200 ответов
type ErrorResponse struct {
	Error APIError `json:"error"`
}
