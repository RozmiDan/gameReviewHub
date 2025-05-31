package entity

import (
	"errors"
	"time"
)

var (
	ErrGameNotFound      = errors.New("game not found")
	ErrGameAlreadyExists = errors.New("game already exists")
	ErrInsertGame        = errors.New("failed to insert game")
	ErrCacheMiss         = errors.New("no data in redis")
)

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

type RequestIDKey struct{}
