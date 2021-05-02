package game

import (
	"github.com/gorilla/websocket"
	"github.com/yigitsadic/minigame/internal/random_generator"
)

type Player struct {
	Identifier    string          `json:"identifier"`
	ClaimedNumber int             `json:"claimed_number"`
	Conn          *websocket.Conn `json:"-"`
}

func NewPlayer(identifier string, con *websocket.Conn) *Player {
	return &Player{
		Identifier:    identifier,
		ClaimedNumber: random_generator.GenerateRandomNumber(),
		Conn:          con,
	}
}
