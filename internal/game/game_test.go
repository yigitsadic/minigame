package game

import (
	"fmt"
	"github.com/yigitsadic/minigame/graph/model"
	"github.com/yigitsadic/minigame/internal"
	"reflect"
	"testing"
)

func TestGame_PrizeDoubled(t *testing.T) {
	p1 := NewPlayer("ABCD1")
	p2 := NewPlayer("ABCD2")

	// Use buffered channel.
	p1C := make(chan *model.Message, 1)
	p2C := make(chan *model.Message, 1)

	p1.MessageChan = p1C
	p2.MessageChan = p2C

	g := NewGame()

	g.Players[p1.Identifier] = p1
	g.Players[p2.Identifier] = p2

	var a *model.Message
	var b *model.Message

	g.PrizeDoubled()

	a = <-p1.MessageChan
	b = <-p2.MessageChan

	if a == nil {
		t.Errorf("expected to get a message from p1 message channel")
	}

	if b == nil {
		t.Errorf("expected to get a message from p2 message channel")
	}

	if a != nil && a.Text != PrizeDoubledMessage {
		t.Errorf("expected message was %q got=%q", PrizeDoubledMessage, a.Text)
	}

	if b != nil && b.Text != PrizeDoubledMessage {
		t.Errorf("expected message was %q got=%q", PrizeDoubledMessage, b.Text)
	}
}

func TestGame_JoinPlayer(t *testing.T) {
	t.Run("should give an error if room is full", func(a *testing.T) {
		g := NewGame()

		for x := 1; x <= internal.PlayerLimit; x++ {
			p := NewPlayer(fmt.Sprintf("Player %d", x))
			g.Players[p.Identifier] = p
		}

		if len(g.Players) != internal.PlayerLimit {
			a.Errorf("expected player count was %d, but got=%d", internal.PlayerLimit, len(g.Players))
		}

		_, err := g.JoinPlayer("Unlucky Player")
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

		got, err := g.JoinPlayer("Lucky Player")
		if err != nil {
			a.Errorf("unexpected to get an error but got=%s", err)
		}

		if reflect.TypeOf(got).String() != "chan *model.Message" {
			a.Errorf("unexpected return type. expected to get %s", reflect.TypeOf(got))
		}
	})
}
