package model

type Message struct {
	ID            string      `json:"id"`
	Text          string      `json:"text"`
	MessageType   MessageType `json:"messageType"`
	ClaimedNumber *int        `json:"claimedNumber"`
	PrizeWon      *int        `json:"prizeWon"`
}

type MessageType string

const (
	MessageTypeInitial     MessageType = "INITIAL"
	MessageTypeUserJoined  MessageType = "USER_JOINED"
	MessageTypeWinnerFound MessageType = "WINNER_FOUND"
	MessageTypeDoublePrize MessageType = "DOUBLE_PRIZE"
	MessageTypeGameEnded   MessageType = "GAME_ENDED"
)
