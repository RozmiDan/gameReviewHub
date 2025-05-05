package entity

import "errors"

var (
	ErrInvalidUUID        = errors.New("entered uuid is invalid")
	ErrInternalRating     = errors.New("rating service error")
	ErrServiceUnavailable = errors.New("rating service unavailable")
	ErrBrokerUnavailable  = errors.New("broker service unavailable")
)

type GameRating struct {
	GameID        string  `json:"gameid"`
	AverageRating float64 `json:"average_rating"`
	RatingsCount  int64   `json:"ratings_count"`
}

type RatingMessage struct {
	GameID string `json:"game_id"`
	UserID string `json:"user_id"`
	Rating int32  `json:"rating"`
}
