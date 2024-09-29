package game

import (
	"github.com/create-go-app/fiber-go-template/app/dao"
	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/pkg/utils"
	"math/rand"
	"strconv"
	"strings"
)

type Service interface {
	dao.Dao
	Value(cards []string) int
	Deal(game *models.Game)
	Hit(game *models.Game)
	Stay(game *models.Game)
	CreateDeck(game *models.Game)
}

type service struct {
	dao.Dao
}

func NewService() (Service, error) {
	gameDao, err := dao.NewDao()
	if err != nil {
		return nil, err
	}
	return &service{gameDao}, nil
}

var (
	suits = []string{
		"\u2660",
		"\u2663",
		"\u2665",
		"\u2666",
	}

	faces = []string{
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"10",
		"A",
		"J",
		"Q",
		"K",
	}
)

func (s *service) CreateDeck(game *models.Game) {
	for _, suit := range suits {
		for _, face := range faces {
			game.Deck = append(game.Deck, face+suit)
		}
	}
	for i := 0; i < len(game.Deck); i++ {
		j := rand.Intn(len(game.Deck))
		origCard := game.Deck[i]
		game.Deck[i] = game.Deck[j]
		game.Deck[j] = origCard
	}
}

func (s *service) Stay(game *models.Game) {
	for s.Value(game.PlayerCards) < 17 {
		game.DealerCards = append(game.DealerCards, utils.Pop((*[]string)(&game.Deck)))
	}
}

func (s *service) Hit(game *models.Game) {
	game.PlayerCards = append(game.PlayerCards, utils.Pop((*[]string)(&game.Deck)))
}

func (s *service) Deal(game *models.Game) {
	game.PlayerCards = append(game.PlayerCards, utils.Pop((*[]string)(&game.Deck)))
	game.DealerCards = append(game.DealerCards, utils.Pop((*[]string)(&game.Deck)))
	game.PlayerCards = append(game.PlayerCards, utils.Pop((*[]string)(&game.Deck)))
	game.DealerCards = append(game.DealerCards, utils.Pop((*[]string)(&game.Deck)))
}

func (s *service) Value(cards []string) int {
	var retval = 0
	var hasAce = false
	for _, card := range cards {
		intVal, err := strconv.Atoi(card)
		if err == nil {
			retval += intVal
			continue
		}
		if strings.Contains(card, "J") || strings.Contains(card, "Q") || strings.Contains(card, "K") {
			retval += 10
			continue
		}
		if strings.Contains(card, "A") {
			hasAce = true
		}
	}
	if hasAce {
		for _, card := range cards {
			if strings.Contains(card, "A") {
				if retval+11 > 21 {
					retval += 1
				} else {
					retval += 11
				}
			}
		}
	}
	return retval
}
