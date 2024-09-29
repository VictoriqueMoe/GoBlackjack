package controllers

import (
	"github.com/create-go-app/fiber-go-template/app/game"
	"github.com/create-go-app/fiber-go-template/app/models"
)

type Service struct {
	MainGameService game.MainGameService
}

func NewService(gameService game.MainGameService) *Service {
	return &Service{
		MainGameService: gameService,
	}
}

func (s *Service) GetAllRoutes() []models.FSetupRoute {
	all := []models.FSetupRoute{}
	all = append(all, s.getAllGameRoutes()...)

	return all
}
