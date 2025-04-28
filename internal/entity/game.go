package entity

import "time"

type Game struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Genre       string     `json:"genre"`
	Creator     string     `json:"creator"`
	Description string     `json:"description"`
	Rating      GameRating `json:"rating"`
	ReleaseDate time.Time  `json:"releasedate"`
}

type GameInList struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Genre  string  `json:"genre"`
	Rating float64 `json:"rating"`
}
