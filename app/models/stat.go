package models

type Stat struct {
	Device string `gorm:"primary_key;type:text;unique;not null"`
	Wins   int    `gorm:"type:integer"`
	Loses  int    `gorm:"type:integer"`
	Draws  int    `gorm:"type:integer"`
}

func NewStat(device string, wins int, loses int, draws int) *Stat {
	return &Stat{
		Device: device,
		Wins:   wins,
		Loses:  loses,
		Draws:  draws,
	}
}

func (Stat) TableName() string {
	return "Stat"
}
