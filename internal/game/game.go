package game

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/yigitsadic/minigame/graph/model"
	"github.com/yigitsadic/minigame/internal"
	"github.com/yigitsadic/minigame/internal/random_generator"
	"sync"
	"time"
)

var (
	InvalidGameError             = errors.New("invalid game")
	UserLimitReachedError        = errors.New("maximum user limit reached")
	InvalidPlayerIdentifierError = errors.New("invalid player identifier")
)

type Game struct {
	Id        string
	CreatedAt time.Time

	Players        map[string]int
	PlayerChannels map[string]chan *model.Message
	WinnerNumber   int

	CurrentPrize int

	LastWinnerCheck time.Time
	NextWinnerCheck time.Time

	Mu sync.Mutex
}

func (g *Game) PrizeDoubled() {
	for _, c := range g.PlayerChannels {
		c <- &model.Message{
			ID:          uuid.NewString(),
			Text:        "Prize is doubled. Wish you luck.",
			MessageType: model.MessageTypeDoublePrize,
		}
	}
}

func (g *Game) JoinPlayer(gameId, identifier string) (chan *model.Message, error) {
	if gameId != g.Id {
		return nil, InvalidGameError
	}

	if len(g.Players) >= internal.PlayerLimit {
		return nil, UserLimitReachedError
	}

	g.Mu.Lock()
	g.Players[identifier] = random_generator.GenerateRandomNumber()
	g.PlayerChannels[identifier] = make(chan *model.Message)
	g.Mu.Unlock()

	return g.PlayerChannels[identifier], nil
}

func (g *Game) PublishClaimedNumber(identifier string) error {
	c, ok1 := g.PlayerChannels[identifier]
	n, ok2 := g.Players[identifier]

	if ok1 && ok2 {
		c <- &model.Message{
			ID:            uuid.NewString(),
			Text:          fmt.Sprintf("Your number is %d. You will notified if you win.", n),
			MessageType:   model.MessageTypeInitial,
			ClaimedNumber: &n,
		}

		return nil
	} else {
		return InvalidPlayerIdentifierError
	}
}
