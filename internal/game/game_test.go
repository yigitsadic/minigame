package game

import (
	"github.com/yigitsadic/minigame/graph/model"
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
