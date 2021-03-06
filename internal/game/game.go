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
	UserLimitReachedError = errors.New("maximum user limit reached")
)

type Game struct {
	Id        string
	CreatedAt time.Time

	Players      map[string]*Player
	WinnerNumber int

	Done   chan bool
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

		Done:         make(chan bool, 1),
		Winner:       make(chan *Player, 1),
		WinnerNumber: random_generator.GenerateRandomNumber(),

		CurrentPrize: internal.StartingPrize,

		LastWinnerCheck: time.Now(),
		NextWinnerCheck: time.Now().Add(time.Minute * internal.TryMinute),
		// TickerChan:      time.Tick(time.Minute * internal.TryMinute),
		TickerChan: time.Tick(time.Second * 10),
	}
}

// Sends an event to inform that prize is doubled.
func (g *Game) PrizeDoubled() {
	g.CurrentPrize *= 2

	event := &Event{
		EType:   EventPrizeDoubled,
		Payload: &PrizeDoubledPayload{NewPrize: g.CurrentPrize},
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
// Handles game. Waits for game finished or winner found signal.
func (g *Game) HandleGame() {

	for {
		if g.CurrentStep >= internal.TryCount {
			g.Done <- true
		}

		select {
		case <-g.Done:
			for _, p := range g.Players {
				evt := &Event{
					EType:   EventGameStopped,
					Player:  p,
					Payload: nil,
				}

				p.Conn.WriteJSON(evt)
			}

			return
		case <-g.TickerChan:
			g.CurrentStep++

			go g.HandleWinner()
			go g.PrizeDoubled()
		case winner := <-g.Winner:
			evt := &Event{
				EType:  EventWinnerFound,
				Player: winner,
				Payload: &WinnerFoundPayload{
					ClaimedPrize: g.CurrentPrize,
				},
			}

			g.Events <- evt
			g.Done <- true
		case evt := <-g.Events:
			if evt.Player != nil {
				evt.Player.Conn.WriteJSON(evt)
			} else {
				for _, p := range g.Players {
					p.Conn.WriteJSON(evt)
				}
			}
		}
	}
}

// Returns winning player if exists.
func (g *Game) WinningPlayer() *Player {
	for _, p := range g.Players {
		if g.WinnerNumber == p.ClaimedNumber {
			return p
		}
	}

	return nil
}

// If winner found, publishes event.
func (g *Game) HandleWinner() {
	winner := g.WinningPlayer()

	if winner == nil {
		return
	}

	g.Mu.Lock()
	defer g.Mu.Unlock()

	go func() {
		g.Winner <- winner
	}()

	go func() {
		evt := &Event{
			EType:   EventWinnerFound,
			Player:  winner,
			Payload: &WinnerFoundPayload{ClaimedPrize: g.CurrentPrize},
		}

		g.Events <- evt
	}()
}
