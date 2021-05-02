package game

type EventType int

const (
	EventPlayerJoined EventType = iota
	EventPrizeDoubled
	EventWinnerFound
	EventGameStopped
)

type Event struct {
	EType   EventType   `json:"event_type"`
	Player  *Player     `json:"player"`
	Payload interface{} `json:"payload"`
}

type PrizeDoubledPayload struct {
	NewPrize int `json:"new_prize"`
}

type PlayerJoinedPayload struct {
	ClaimedNumber int `json:"claimed_number"`
	CurrentPrize  int `json:"current_prize"`
}

type WinnerFoundPayload struct {
	ClaimedPrize int `json:"claimed_prize"`
}
