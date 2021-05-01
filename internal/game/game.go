package game

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/yigitsadic/minigame/graph/model"
	"github.com/yigitsadic/minigame/internal"
	"github.com/yigitsadic/minigame/internal/random_generator"
	"log"
	"sync"
	"time"
)

var (
	UserLimitReachedError        = errors.New("maximum user limit reached")
	InvalidPlayerIdentifierError = errors.New("invalid player identifier")
)

const (
	PrizeDoubledMessage = "Prize is doubled. Wish you luck."
)

type Game struct {
	Id        string
	CreatedAt time.Time

	Players      map[string]*Player
	WinnerNumber int

	CurrentStep  int
	CurrentPrize int

	LastWinnerCheck time.Time
	NextWinnerCheck time.Time

	Mu         sync.Mutex
	TickerChan <-chan time.Time
}

// Initializes a new game with default options.
func NewGame() *Game {
	return &Game{
		Id:           uuid.NewString(),
		CreatedAt:    time.Now(),
		Players:      make(map[string]*Player),
		WinnerNumber: random_generator.GenerateRandomNumber(),

		CurrentPrize: internal.StartingPrize,

		LastWinnerCheck: time.Now(),
		NextWinnerCheck: time.Now().Add(time.Minute * internal.TryMinute),
		TickerChan:      time.Tick(time.Minute * internal.TryMinute),
	}
}

// Sends a message to all participants that prize is doubled.
func (g *Game) PrizeDoubled() {
	for _, player := range g.Players {
		go func(p *Player) {
			defer func() {
				recover()
			}()

			p.MessageChan <- &model.Message{
				ID:          uuid.NewString(),
				Text:        PrizeDoubledMessage,
				MessageType: model.MessageTypeDoublePrize,
			}
		}(player)
	}
}

// Joins a player with given identifier if there is a room for him/her.
func (g *Game) JoinPlayer(identifier string) (chan *model.Message, error) {
	if len(g.Players) >= internal.PlayerLimit {
		return nil, UserLimitReachedError
	}

	g.Mu.Lock()

	p := &Player{
		Identifier:    identifier,
		ClaimedNumber: random_generator.GenerateRandomNumber(),
		MessageChan:   make(chan *model.Message),
	}

	g.Players[identifier] = p
	g.Mu.Unlock()

	return p.MessageChan, nil
}

func (g *Game) PublishClaimedNumber(identifier string) error {
	p, ok := g.Players[identifier]

	if ok {
		p.MessageChan <- &model.Message{
			ID:            uuid.NewString(),
			Text:          fmt.Sprintf("Your number is %d. You will notified if you win.", p.ClaimedNumber),
			MessageType:   model.MessageTypeInitial,
			ClaimedNumber: &p.ClaimedNumber,
		}

		return nil
	} else {
		return InvalidPlayerIdentifierError
	}
}

func (g *Game) HandleGameTicker() {
	for t := range g.TickerChan {
		if g.CurrentStep >= internal.TryCount {
			log.Println("Game stopped.")
			break
		}

		g.CurrentStep++

		g.Mu.Lock()

		if p := g.WinningPlayer(); p != nil {
			log.Println("Winner found")
		}

		g.LastWinnerCheck = t
		g.NextWinnerCheck = t.Add(time.Minute * internal.TryMinute)
		g.CurrentPrize *= 2

		g.Mu.Unlock()

		g.PrizeDoubled()
	}
}

// Returns winning player if exists.
func (g *Game) WinningPlayer() *Player {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	for _, p := range g.Players {
		if g.WinnerNumber == p.ClaimedNumber {
			return p
		}
	}

	return nil
}

func (g *Game) PublishToWinner() {
	panic("implement me")
}

func (g *Game) PublishToLosers() {
	panic("implement me")
}

func (g *Game) CloseAllChannels() {
	for _, p := range g.Players {
		close(p.MessageChan)
	}
}
