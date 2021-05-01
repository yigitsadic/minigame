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
	MessageTypeYouWin      MessageType = "YOU_WIN"
	MessageTypeDoublePrize MessageType = "DOUBLE_PRIZE"
	MessageTypeGameEnded   MessageType = "GAME_ENDED"
)
