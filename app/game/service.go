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
	Stay(game *models.Game) (playerValue int, dealerValue int, err error)
	CreateDeck(game *models.Game)
	GetActiveGame(deviceId string) (*models.Game, error)
	GetGameFromToken(deviceId string) (*models.Game, error)
	GetAllGames(deviceId string) ([]models.Game, error)
}

type service struct {
	dao         dao.Dao
	statService stats.StatService
}

func NewService(
	statService stats.StatService,
	dao dao.Dao,
) (MainGameService, error) {
	return &service{
		dao:         dao,
		statService: statService,
	}, nil
}

func (s *service) GetAllGames(deviceId string) ([]models.Game, error) {
	return s.dao.GetAllGames(deviceId)
}

func (s *service) GetGameFromToken(token string) (*models.Game, error) {
	return s.dao.RetrieveGame(token)
}

func (s *service) GetActiveGame(deviceId string) (*models.Game, error) {
	return s.dao.RetrieveActiveGame(deviceId)
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

func (s *service) Stay(game *models.Game) (playerValue int, dealerValue int, err error) {
	for s.CalculateScore(game.DealerCards) < 17 {
		game.DealerCards = append(game.DealerCards, utils.Pop(&game.Deck))
	}
	playerValue = s.CalculateScore(game.PlayerCards)
	dealerValue = s.CalculateScore(game.DealerCards)

	if dealerValue > 21 {
		game.Status = models.DealerBust
		log.Info("DEALER BUST")
		_, err = s.statService.UpdateStats(game.Device, stats.Win)
	} else if playerValue > dealerValue {
		game.Status = models.PlayerWins
		log.Info("WIN")
		_, err = s.statService.UpdateStats(game.Device, stats.Win)
	} else if dealerValue > playerValue {
		game.Status = models.DealerWins
		log.Info("LOSE")
		_, err = s.statService.UpdateStats(game.Device, stats.Loss)
	} else {
		game.Status = models.Draw
		log.Info("DRAW")
		_, err = s.statService.UpdateStats(game.Device, stats.Draw)
	}

	_, err = s.dao.SaveOrUpdateGame(*game)
	return
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
