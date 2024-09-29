package controllers

import (
	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (s *Service) getAllGameRoutes() []models.FSetupRoute {
	return []models.FSetupRoute{
		s.SetupDealRoute,
		s.setupHitRoute,
	}
}

// setupHitRoute.
// @Description Draw a new card from the dealer
// @Summary hit move
// @Tags Game
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseMsg
// @Success 500 {object} models.ErrorMsg
// @Success 400 {object} models.ErrorMsg
// @Router /api/v1/hit [get]
func (s *Service) setupHitRoute(routeGroup fiber.Router) {
	routeGroup.Get("/hit", s.hit)
}

func (s *Service) hit(c *fiber.Ctx) error {
	ip := c.IP()
	ua := c.Get("User-Agent")
	deviceId := utils.DeviceHash(ip, ua)
	retGame, err := s.MainGameService.GetGame(deviceId, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorMsg(
			"internal server error",
			err,
			fiber.StatusInternalServerError,
		))
	}

	if retGame == nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorMsg(
			"No game is in progress",
			err,
			fiber.StatusBadRequest,
		))
	}

	err = s.MainGameService.Hit(retGame)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorMsg(
			"Hit failed",
			err,
			fiber.StatusInternalServerError,
		))
	}

	resp := models.NewResponseMsg(
		retGame.Device,
		retGame.PlayerCards,
		[]string{},
		s.MainGameService.CalculateScore(retGame.PlayerCards),
		0,
		retGame.Status,
	)

	return c.Status(fiber.StatusOK).JSON(resp)

}

// SetupDealRoute Deal deals a new move.
// @Description Start a new game or deal for an existing one.
// @Summary start game
// @Tags Game
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseMsg
// @Success 500 {object} models.ErrorMsg
// @Router /api/v1/deal [get]
func (s *Service) SetupDealRoute(routeGroup fiber.Router) {
	routeGroup.Get("/deal", s.deal)
}

func (s *Service) deal(c *fiber.Ctx) error {
	ip := c.IP()
	ua := c.Get("User-Agent")
	deviceId := utils.DeviceHash(ip, ua)
	retGame, err := s.MainGameService.GetGame(deviceId, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorMsg(
			"internal server error",
			err,
			fiber.StatusInternalServerError,
		))
	}

	if retGame == nil {
		newGame, err := s.MainGameService.NewGame(deviceId)
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
			s.MainGameService.CalculateScore(retGame.PlayerCards),
			0,
			retGame.Status,
		))
}
