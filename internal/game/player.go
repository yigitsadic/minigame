package game

import (
	"github.com/yigitsadic/minigame/graph/model"
	"github.com/yigitsadic/minigame/internal/random_generator"
)

type Player struct {
	Identifier    string
	ClaimedNumber int
	MessageChan   chan *model.Message
}

func NewPlayer(identifier string) *Player {
	return &Player{
		Identifier:    identifier,
		ClaimedNumber: random_generator.GenerateRandomNumber(),
		MessageChan:   make(chan *model.Message),
	}
}
