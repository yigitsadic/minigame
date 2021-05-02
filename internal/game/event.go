package game

type EventType int

const (
	EventPlayerJoined EventType = iota
	EventPrizeDoubled
	EventWinnerFound
	EventGameStopped
)

type Event struct {
	EType   EventType
	Player  *Player
	Payload interface{}
}

type PrizeDoubledPayload struct {
	NewPrize int
}

type PlayerJoinedPayload struct {
	ClaimedNumber int
	CurrentPrize  int
}
