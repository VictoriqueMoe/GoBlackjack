package models

import "github.com/google/uuid"

type Game struct {
	Token       uuid.UUID  `gorm:"type:TEXT;primary_key;unique"`
	Device      string     `gorm:"type:TEXT"`
	Status      PlayStatus `gorm:"type:TEXT"`
	StartedOn   int64      `gorm:"type:integer"`
	Deck        []string   `gorm:"type:TEXT;serializer:json"`
	DealerCards []string   `gorm:"type:TEXT;serializer:json"`
	PlayerCards []string   `gorm:"type:TEXT;serializer:json"`
}

func NewGame(
	device string,
	status PlayStatus,
	startedOn int64,
) *Game {
	return newGame(
		uuid.New(),
		device,
		status,
		startedOn,
		[]string{},
		[]string{},
		[]string{},
	)
}

func newGame(
	token uuid.UUID,
	device string,
	status PlayStatus,
	startedOn int64,
	deck []string,
	dealerCards []string,
	playerCards []string,
) *Game {
	return &Game{
		Token:       token,
		Device:      device,
		Status:      status,
		StartedOn:   startedOn,
		Deck:        deck,
		DealerCards: dealerCards,
		PlayerCards: playerCards,
	}
}

func (Game) TableName() string {
	return "Game"
}
