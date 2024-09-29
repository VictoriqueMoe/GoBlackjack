package dao

import (
	"errors"
	"github.com/create-go-app/fiber-go-template/app/models"
	"gorm.io/gorm"
)

type StatDao interface {
	GetStat(deviceId string, tx ...*gorm.DB) (*models.Stat, error)
	SaveStat(stat models.Stat, tx ...*gorm.DB) (models.Stat, error)
}

func (d dao) GetStat(deviceId string, tx ...*gorm.DB) (*models.Stat, error) {
	var stat models.Stat
	err := d.getDb(tx...).First(&stat, "device = ?", deviceId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	return &stat, nil
}

func (d dao) SaveStat(stat models.Stat, tx ...*gorm.DB) (models.Stat, error) {
	err := d.getDb(tx...).Save(&stat).Error
	return stat, err
}
