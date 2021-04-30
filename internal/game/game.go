package game

import (
	"github.com/yigitsadic/minigame/graph/model"
	"sync"
	"time"
)

type Game struct {
	Id        string
	CreatedAt time.Time

	Players        map[string]int
	PlayerChannels map[string]chan *model.Message
	WinnerNumber   int

	CurrentPrize    int
	NextWinnerCheck time.Time

	Mu sync.Mutex
}
