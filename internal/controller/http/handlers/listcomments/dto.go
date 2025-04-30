package handlers

import "github.com/RozmiDan/gameReviewHub/internal/entity"

type Pagination struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
	Count  int   `json:"count,omitempty"`
	Total  int   `json:"total,omitempty"`
}

// ListGamesResponse — обёртка для GET /games
type ListCommentsResponse struct {
	Data []entity.Comment `json:"data"`
	Meta *Pagination      `json:"meta,omitempty"`
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
