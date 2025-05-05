package handlers

type CreateGameRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Genre       string `json:"genre"`
	Creator     string `json:"creator"`
	Description string `json:"description"`
	ReleaseDate string `json:"release_date"` // or time.Time + правильный UnmarshalJSON
}

type CreateGameResponse struct {
	ID string `json:"id"`
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
