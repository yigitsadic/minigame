package model

type Game struct {
	ID              string `json:"id"`
	CurrentPrize    int    `json:"currentPrize"`
	NextWinnerCheck string `json:"nextWinnerCheck"`
	LastWinnerCheck string `json:"lastWinnerCheck"`
	CreatedAt       string `json:"createdAt"`
}
