package entity

import "errors"

var (
	ErrGameNotFound = errors.New("game not found")
	ErrInvalidUUID  = errors.New("entered uuid is invalid")
)

type GameRating struct {
	GameID        string  `json:"gameid"`
	AverageRating float64 `json:"average_rating"`
	RatingsCount  int64   `json:"ratings_count"`
}
