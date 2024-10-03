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
		s.setupStayRoute,
		s.setupStatRoute,
		s.setupHistoryRoute,
	}
}

// setupHistoryRoute.
// @Description Get all the games you played
// @Summary get all games played
// @Tags Game
// @Accept json
// @Produce json
// @Success 200 {object} []models.StatusMsg
// @Success 500 {object} models.ErrorMsg
// @Success 404 {object} models.ErrorMsg
// @Router /api/v1/history [get]
func (s *Service) setupHistoryRoute(routeGroup fiber.Router) {
	routeGroup.Get("/history", s.history)
}

// setupStatRoute.
// @Description Get player stats of wins, loses and draws
// @Summary Get player stats of wins, loses and draws
// @Tags Game
// @Accept json
// @Produce json
// @Success 200 {object} []models.ResponseMsg
// @Success 500 {object} models.ErrorMsg
// @Success 404 {object} models.ErrorMsg
// @Router /api/v1/stats [get]
func (s *Service) setupStatRoute(routeGroup fiber.Router) {
	routeGroup.Get("/stats", s.history)
}

func (s *Service) history(c *fiber.Ctx) error {
	deviceId := getDeviceId(c)
	games, err := s.MainGameService.GetAllGames(deviceId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorMsg(
			"internal server error",
			err,
			fiber.StatusInternalServerError,
		))
	}

	var resp []models.ResponseMsg

	for _, game := range games {
		resp = append(resp, *models.NewResponseMsg(
			game.Token,
			game.Device,
			game.PlayerCards,
			game.DealerCards,
			s.MainGameService.CalculateScore(game.PlayerCards),
			s.MainGameService.CalculateScore(game.DealerCards),
			game.Status,
		))
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (s *Service) stats(c *fiber.Ctx) error {
	deviceId := getDeviceId(c)
	stats, err := s.StatsService.GetStats(deviceId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorMsg(
			"internal server error",
			err,
			fiber.StatusInternalServerError,
		))
	}

	if stats == nil {
		log.Warnf("No stats found for device %s", deviceId)
		return c.Status(fiber.StatusNotFound).JSON(models.NewErrorMsg(
			"No stats found for device",
			nil,
			fiber.StatusNotFound,
		))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewStatusMsg(stats.Wins, stats.Loses, stats.Draws))
}

// setupStayRoute.
// @Description Stop drawing cards and allow dealer to draw cards
// @Summary Stop drawing cards and allow dealer to draw cards
// @Tags Game
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseMsg
// @Success 500 {object} models.ErrorMsg
// @Success 400 {object} models.ErrorMsg
// @Param	token	query	string	false	"token id for a game"
// @Router /api/v1/stay [get]
func (s *Service) setupStayRoute(routeGroup fiber.Router) {
	routeGroup.Get("/stay", s.stay)
}

func (s *Service) stay(c *fiber.Ctx) error {
	token := c.Query("token")
	var retGame *models.Game
	if token != "" {
		game, err := s.MainGameService.GetGameFromToken(token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorMsg(
				"internal server error",
				err,
				fiber.StatusInternalServerError,
			))
		}
		retGame = game
	} else {
		deviceId := getDeviceId(c)
		game, err := s.MainGameService.GetActiveGame(deviceId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorMsg(
				"internal server error",
				err,
				fiber.StatusInternalServerError,
			))
		}
		retGame = game
	}

	if retGame == nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorMsg(
			"No game is in progress",
			nil,
			fiber.StatusBadRequest,
		))
	}

	playerVal, dealerVal, err := s.MainGameService.Stay(retGame)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.NewResponseMsg(
		retGame.Token,
		retGame.Device,
		retGame.PlayerCards,
		retGame.DealerCards,
		playerVal,
		dealerVal,
		retGame.Status,
	))
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
// @Param	token	query	string	false	"token id for a game"
// @Router /api/v1/hit [get]
func (s *Service) setupHitRoute(routeGroup fiber.Router) {
	routeGroup.Get("/hit", s.hit)
}

func (s *Service) hit(c *fiber.Ctx) error {
	deviceId := getDeviceId(c)
	retGame, err := s.MainGameService.GetActiveGame(deviceId)
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
			nil,
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
		retGame.Token,
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
	deviceId := getDeviceId(c)
	retGame, err := s.MainGameService.GetActiveGame(deviceId)
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
			retGame.Token,
			retGame.Device,
			retGame.PlayerCards,
			[]string{},
			s.MainGameService.CalculateScore(retGame.PlayerCards),
			0,
			retGame.Status,
		))
}

func getDeviceId(context *fiber.Ctx) string {
	ip := context.IP()
	ua := context.Get("User-Agent")
	return utils.DeviceHash(ip, ua)
}
