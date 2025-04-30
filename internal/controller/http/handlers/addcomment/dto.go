package handlers

// PostRatingRequest — тело запроса для POST /games/{game_id}/rating
type PostCommentRequest struct {
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}

// APIError — единая структура описания ошибки
type APIError struct {
	Code    string `json:"code"`    // машинно-читаемый код ошибки
	Message string `json:"message"` // человеко-читаемое сообщение
}

// ErrorResponse — обёртка над APIError
type ErrorResponse struct {
	Error APIError `json:"error"`
}
