package controllers

import (
	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

func (s *Service) getAllBookRoutes() []models.FSetupRoute {
	return []models.FSetupRoute{
		s.setupCreateBookRoute,
	}
}

// Deal deals a new move.
// @Description Start a new game or deal for an existing one.
// @Summary start game
// @Tags Game
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseMsg
// @Success 500 {object} models.ErrorMsg
// @Router /v1/deal [get]
func (s *Service) setupCreateBookRoute(routeGroup fiber.Router) {
	routeGroup.Get("/deal", s.deal)
}

func (s *Service) deal(c *fiber.Ctx) error {
	ip := c.IP()
	ua := c.Get("User-Agent")
	deviceId := utils.DeviceHash(ip, ua)
	retGame, err := s.RetrieveGame(deviceId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorMsg(
			"internal server error",
			err,
			fiber.StatusInternalServerError,
		))
	}

	if retGame == nil {
		newGame := *models.NewGame(
			deviceId,
			models.Playing,
			time.Now().Unix(),
		)
		s.CreateDeck(&newGame)
		s.Deal(&newGame)
		newGame, err := s.SaveOrUpdateGame(newGame)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorMsg(
				"internal server error",
				err,
				fiber.StatusInternalServerError,
			))
		}
		retGame = &newGame
	}

	log.Infof("DEAL: %s", retGame.Device)
	return c.Status(fiber.StatusOK).
		JSON(models.NewResponseMsg(
			retGame.Device,
			retGame.PlayerCards,
			[]string{},
			s.Value(retGame.PlayerCards),
			0,
			retGame.Status,
		))
}
