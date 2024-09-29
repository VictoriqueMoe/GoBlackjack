package controllers

import (
	"github.com/create-go-app/fiber-go-template/app/game"
	"github.com/create-go-app/fiber-go-template/app/models"
)

type Service struct {
	game.Service
}

func NewService(local game.Service) *Service {
	return &Service{
		Service: local,
	}
}

func (s *Service) GetAllRoutes() []models.FSetupRoute {
	all := []models.FSetupRoute{}
	all = append(all, s.getAllBookRoutes()...)

	return all
}
