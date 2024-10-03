package dao

import (
	"errors"
	"github.com/create-go-app/fiber-go-template/app/models"
	"gorm.io/gorm"
)

type GameDao interface {
	SaveOrUpdateGame(game models.Game, tx ...*gorm.DB) (models.Game, error)
	RemoveGame(game models.Game, tx ...*gorm.DB) error
	RetrieveGame(token string, tx ...*gorm.DB) (*models.Game, error)
	RetrieveActiveGame(deviceId string, tx ...*gorm.DB) (*models.Game, error)
	GetAllGames(deviceId string, tx ...*gorm.DB) ([]models.Game, error)
}

func (d dao) SaveOrUpdateGame(game models.Game, tx ...*gorm.DB) (models.Game, error) {
	err := d.getDb(tx...).Save(&game).Error
	return game, err
}

func (d dao) GetAllGames(deviceId string, tx ...*gorm.DB) ([]models.Game, error) {
	var games []models.Game
	err := d.getDb(tx...).Where("device = ?", deviceId).Find(&games).Error
	return games, err
}

func (d dao) RemoveGame(game models.Game, tx ...*gorm.DB) error {
	err := d.getDb(tx...).Delete(&game).Error
	return err
}

func (d dao) RetrieveActiveGame(deviceId string, tx ...*gorm.DB) (*models.Game, error) {
	var game models.Game
	err := d.getDb(tx...).
		Where("Status = ?", models.Playing).
		Where("device = ?", deviceId).
		First(&game).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	return &game, err
}

func (d dao) RetrieveGame(token string, tx ...*gorm.DB) (*models.Game, error) {
	var game models.Game
	err := d.getDb(tx...).First(&game, "Token = ?", token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	return &game, err
}
