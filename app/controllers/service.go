package controllers

import (
	"github.com/create-go-app/fiber-go-template/app/dao"
	"github.com/create-go-app/fiber-go-template/app/game"
	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/app/stats"
)

type Service struct {
	MainGameService game.MainGameService
	StatsService    stats.StatService
}

func NewService(dao dao.Dao) *Service {
	statsService, err := stats.NewService(dao)
	if err != nil {
		panic("an error occurred initialising the the stats service")
	}

	gameService, err := game.NewService(statsService, dao)
	if err != nil {
		panic("an error occurred initialising the the game service")
	}
	return &Service{
		MainGameService: gameService,
		StatsService:    statsService,
	}
}

func (s *Service) GetAllRoutes() []models.FSetupRoute {
	all := []models.FSetupRoute{}
	all = append(all, s.getAllGameRoutes()...)

	return all
}
