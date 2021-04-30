package game

import "github.com/yigitsadic/minigame/graph/model"

type Player struct {
	Identifier    string
	ClaimedNumber int
	MessageChan   chan *model.Message
}
