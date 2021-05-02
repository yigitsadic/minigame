package game

import (
	"errors"
	"github.com/google/uuid"
	"github.com/yigitsadic/minigame/internal"
	"github.com/yigitsadic/minigame/internal/random_generator"
	"sync"
	"time"
)

var (
	UserLimitReachedError        = errors.New("maximum user limit reached")
	InvalidPlayerIdentifierError = errors.New("invalid player identifier")
)

const (
	PrizeDoubledMessage = "Prize is doubled. Wish you luck."
	YouWinMessage       = "You have won this game. You are a lucky guy/girl!"
)

type Game struct {
	Id        string
	CreatedAt time.Time

	Players      map[string]*Player
	WinnerNumber int

	Winner chan *Player
	Events chan *Event

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
		Id:        uuid.NewString(),
		CreatedAt: time.Now(),
		Players:   make(map[string]*Player),
		Events:    make(chan *Event),

		WinnerNumber: random_generator.GenerateRandomNumber(),
		Winner:       make(chan *Player, 1),

		CurrentPrize: internal.StartingPrize,

		LastWinnerCheck: time.Now(),
		NextWinnerCheck: time.Now().Add(time.Minute * internal.TryMinute),
		TickerChan:      time.Tick(time.Minute * internal.TryMinute),
	}
}

// Sends an event to inform that prize is doubled.
func (g *Game) PrizeDoubled() {
	event := &Event{
		EType:   EventPrizeDoubled,
		Payload: &PrizeDoubledPayload{NewPrize: g.CurrentPrize * 2},
	}

	g.Events <- event
}

// Joins a player with given identifier if there is a room for him/her.
func (g *Game) JoinPlayer(p *Player) error {
	if len(g.Players) >= internal.PlayerLimit {
		return UserLimitReachedError
	}

	evt := &Event{
		EType:  EventPlayerJoined,
		Player: p,
		Payload: &PlayerJoinedPayload{
			ClaimedNumber: p.ClaimedNumber,
			CurrentPrize:  g.CurrentPrize,
		},
	}

	g.Mu.Lock()

	g.Players[p.Identifier] = p
	g.Events <- evt

	g.Mu.Unlock()

	return nil
}

// TODO: Refactor!
func (g *Game) HandleGameTicker() {
	panic("implement me!")
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

// TODO: Refactor!
// Publishes you win message to winner reading Game.Winner.
func (g *Game) PublishToWinner() {
	/*
		if g.Winner == nil {
			return
		}

		message := &model.Message{
			ID:          uuid.NewString(),
			Text:        YouWinMessage,
			MessageType: model.MessageTypeYouWin,
			PrizeWon:    &g.CurrentPrize,
		}

		go func(m *model.Message) {
			defer func() {
				recover()
			}()

			w := <-g.Winner
			w.MessageChan <- m
		}(message)

	*/
}

// TODO: Refactor!
func (g *Game) PublishToLosers() {
	panic("implement me")
}
