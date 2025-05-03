package entity

import (
	"errors"
	"time"
)

var (
	ErrInsertComment    = errors.New("failed to insert comment")
	ErrUnidentified     = errors.New("unidentified error")
	ErrInternal         = errors.New("internal error")
	ErrInternalComments = errors.New("could not fetch comments")
	ErrTimeout          = errors.New("timeout exceeded")
)

type Comment struct {
	ID        string    `json:"id"`
	GameID    string    `json:"game_id"`
	UserID    string    `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
