package game

import (
	"github.com/create-go-app/fiber-go-template/app/dao"
	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/app/stats"
	"github.com/create-go-app/fiber-go-template/pkg/utils"
	"github.com/gofiber/fiber/v2/log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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

	re = regexp.MustCompile("[0-9]+")
)

type MainGameService interface {
	NewGame(deviceId string) (models.Game, error)
	CalculateScore(cards []string) int
	Deal(game *models.Game)
	Hit(game *models.Game) error
	Stay(game *models.Game)
	CreateDeck(game *models.Game)
	GetGame(deviceId string, active bool) (*models.Game, error)
}

type service struct {
	dao         dao.Dao
	statService stats.StatService
}

func NewService() (MainGameService, error) {
	daoService, err := dao.NewDao()
	if err != nil {
		return nil, err
	}
	statService, err := stats.NewService(daoService)
	if err != nil {
		return nil, err
	}
	return &service{
		dao:         daoService,
		statService: statService,
	}, nil
}

func (s *service) GetGame(deviceId string, active bool) (*models.Game, error) {
	if active {
		return s.dao.RetrieveActiveGame(deviceId)
	}
	return s.dao.RetrieveGame(deviceId)
}

func (s *service) NewGame(deviceId string) (models.Game, error) {
	newGame := *models.NewGame(
		deviceId,
		models.Playing,
		time.Now().Unix(),
	)
	s.CreateDeck(&newGame)
	s.Deal(&newGame)
	return s.dao.SaveOrUpdateGame(newGame)
}

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
	for s.CalculateScore(game.PlayerCards) < 17 {
		game.DealerCards = append(game.DealerCards, utils.Pop(&game.Deck))
	}
}

func (s *service) Hit(game *models.Game) error {
	game.PlayerCards = append(game.PlayerCards, utils.Pop(&game.Deck))

	if s.CalculateScore(game.PlayerCards) > 21 {
		game.Status = models.Bust
		log.Info("BUST")
	}

	if game.Status == models.Bust {
		_, err := s.bust(game)
		if err != nil {
			return err
		}
		return nil
	}

	_, err := s.dao.SaveOrUpdateGame(*game)
	if err != nil {
		return err
	}

	log.Infof("HIT: %s", game.Device)
	return nil
}

func (s *service) bust(game *models.Game) (*models.Game, error) {
	_, err := s.statService.UpdateStats(game.Device, stats.Loss)
	if err != nil {
		return nil, err
	}

	updatedGame, err := s.dao.SaveOrUpdateGame(*game)
	if err != nil {
		return nil, err
	}

	return &updatedGame, nil
}

func (s *service) Deal(game *models.Game) {
	game.PlayerCards = append(game.PlayerCards, utils.Pop(&game.Deck))
	game.DealerCards = append(game.DealerCards, utils.Pop(&game.Deck))
	game.PlayerCards = append(game.PlayerCards, utils.Pop(&game.Deck))
	game.DealerCards = append(game.DealerCards, utils.Pop(&game.Deck))
}

func getNumberFromCard(card string) int {
	if re.MatchString(card) {
		num, _ := strconv.Atoi(re.FindString(card))
		return num
	}
	return 0
}

func (s *service) CalculateScore(cards []string) int {
	var retval = 0
	var hasAce = false
	for _, card := range cards {
		retval += getNumberFromCard(card)

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
