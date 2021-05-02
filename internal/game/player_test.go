package game

import (
	"testing"
)

func TestNewPlayer(t *testing.T) {
	got := NewPlayer("ABCDE", nil)

	if got.Identifier != "ABCDE" {
		t.Errorf("expected to initialize as expected. identifier expected=%s got=%s", "ABCDE", got.Identifier)
	}
}
