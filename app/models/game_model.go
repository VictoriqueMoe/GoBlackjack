package models

type Game struct {
	Device      string      `gorm:"type:TEXT;primary_key;unique"`
	Status      PlayStatus  `gorm:"type:TEXT"`
	StartedOn   int64       `gorm:"type:integer"`
	Deck        StringArray `gorm:"type:TEXT;serializer:json"`
	DealerCards StringArray `gorm:"type:TEXT;serializer:json"`
	PlayerCards StringArray `gorm:"type:TEXT;serializer:json"`
}

func NewGame(
	device string,
	status PlayStatus,
	startedOn int64,
) *Game {
	return newGame(
		device,
		status,
		startedOn,
		[]string{},
		[]string{},
		[]string{},
	)
}

func newGame(
	device string,
	status PlayStatus,
	startedOn int64,
	deck []string,
	dealerCards []string,
	playerCards []string,
) *Game {
	return &Game{
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
