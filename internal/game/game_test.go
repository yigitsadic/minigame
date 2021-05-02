package game

import (
	"fmt"
	"github.com/yigitsadic/minigame/internal"
	"testing"
)

func TestGame_PrizeDoubled(t *testing.T) {
	t.Run("it should publish prize doubled event gracefully", func(c *testing.T) {
		c.Parallel()

		startingPrize := 20
		expectedPrize := startingPrize * 2

		g := NewGame()
		g.CurrentPrize = startingPrize
		g.Events = make(chan *Event, 1)

		g.PrizeDoubled()

		evt := <-g.Events

		if evt.EType != EventPrizeDoubled {
			c.Errorf("expected event type was=%d but got=%d", EventPrizeDoubled, evt.EType)
		}

		if p, ok := evt.Payload.(*PrizeDoubledPayload); ok {
			if p.NewPrize != expectedPrize {
				c.Errorf("prize not doubled. expected=%d but got=%d", expectedPrize, p.NewPrize)
			}
		} else {
			c.Errorf("expected payload not satisfied. payload=%v", evt.Payload)
		}

		if evt.Player != nil {
			c.Errorf("unexpected to see a player.")
		}
	})
}

func TestGame_JoinPlayer(t *testing.T) {
	t.Run("should give an error if room is full", func(a *testing.T) {
		g := NewGame()

		for x := 1; x <= internal.PlayerLimit; x++ {
			p := NewPlayer(fmt.Sprintf("Player %d", x), nil)
			g.Players[p.Identifier] = p
		}

		if len(g.Players) != internal.PlayerLimit {
			a.Errorf("expected player count was %d, but got=%d", internal.PlayerLimit, len(g.Players))
		}

		err := g.JoinPlayer(NewPlayer("Unlucky Player", nil))
		if err == nil {
			a.Errorf("expected to get an error, but got nothing")
		}

		if err != UserLimitReachedError {
			a.Errorf("expected error was %s but got %s", UserLimitReachedError, err)
		}

		a.Parallel()
	})

	t.Run("should join player to room", func(a *testing.T) {
		a.Parallel()

		g := NewGame()
		g.Events = make(chan *Event, 1)

		p := NewPlayer("Lucky Player", nil)

		err := g.JoinPlayer(p)
		if err != nil {
			a.Errorf("unexpected to get an error but got=%s", err)
		}

		evt := <-g.Events

		if evt.EType != EventPlayerJoined {
			a.Errorf("expected event type was=%d but got=%d", EventPlayerJoined, evt.EType)
		}

		if evt.Player.Identifier != p.Identifier {
			a.Errorf("expected to contain player information")
		}

		if p, ok := evt.Payload.(*PlayerJoinedPayload); ok {
			if p.ClaimedNumber != p.ClaimedNumber {
				a.Errorf("expected to contain claimed number in payload")
			}

			if p.CurrentPrize != g.CurrentPrize {
				a.Errorf("expected to contain current prize in payload")
			}
		} else {
			a.Errorf("expected payload not satisfied.")
		}
	})
}

func TestGame_WinningPlayer(t *testing.T) {
	t.Run("it should return winner if exists", func(a *testing.T) {
		a.Parallel()

		g := NewGame()

		p := NewPlayer("ABC", nil)
		g.Players[p.Identifier] = p

		p.ClaimedNumber = g.WinnerNumber

		got := g.WinningPlayer()

		if got != p {
			a.Errorf("expected winner was player but got=%v", got)
		}
	})

	t.Run("it should return nil if no winner found", func(a *testing.T) {
		a.Parallel()

		g := NewGame()

		p := NewPlayer("ABC", nil)
		g.Players[p.Identifier] = p

		p.ClaimedNumber = g.WinnerNumber + 1

		got := g.WinningPlayer()

		if got != nil {
			a.Errorf("expected no winner but got=%v", got)
		}
	})
}

func TestGame_HandleWinner(t *testing.T) {
	t.Run("it should send event if winner found", func(a *testing.T) {
		a.Parallel()

		g := NewGame()
		g.Winner = make(chan *Player, 1)
		g.Events = make(chan *Event, 1)

		p := NewPlayer("ABC", nil)

		p.ClaimedNumber = g.WinnerNumber

		g.Players[p.Identifier] = p

		g.HandleWinner()

		evt := <-g.Events
		winner := <-g.Winner

		if evt.EType != EventWinnerFound {
			a.Errorf("expected to get winner found event type")
		}

		if evt.Player.Identifier != p.Identifier {
			a.Errorf("expected to see correct winner")
		}

		if payload, ok := evt.Payload.(*WinnerFoundPayload); ok {
			if payload.ClaimedPrize != g.CurrentPrize {
				a.Errorf("expected prize not satisfied")
			}
		} else {
			a.Errorf("expected payload not satisfied")
		}

		if winner.Identifier != p.Identifier {
			a.Errorf("expected to see correct player as winner")
		}
	})

	t.Run("it should handle gracefully if no winner found", func(a *testing.T) {
		a.Parallel()

		g := NewGame()
		g.HandleWinner()
	})
}
