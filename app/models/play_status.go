package models

type PlayStatus string

const (
	DealerBust PlayStatus = "Dealer Bust"
	PlayerWins PlayStatus = "Player Wins"
	DealerWins PlayStatus = "Dealer Wins"
	Draw       PlayStatus = "Draw"
	Playing    PlayStatus = "Playing"
)
