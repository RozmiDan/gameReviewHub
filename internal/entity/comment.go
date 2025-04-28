package entity

import "time"

type Comment struct {
	ID        string    `json:"id"`
	GameID    string    `json:"game_id"`
	UserID    string    `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
