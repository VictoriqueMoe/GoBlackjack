package stats

import (
	"github.com/create-go-app/fiber-go-template/app/dao"
	"github.com/create-go-app/fiber-go-template/app/models"
)

type StatService interface {
	UpdateStats(deviceId string, action Action) (*models.Stat, error)
}

type statService struct {
	dao dao.Dao
}

func NewService(daoService dao.Dao) (StatService, error) {
	return &statService{
		daoService,
	}, nil
}

func (s *statService) UpdateStats(deviceId string, action Action) (*models.Stat, error) {
	stat, err := s.dao.GetStat(deviceId)
	if err != nil {
		return nil, err
	}

	if stat == nil {
		stat = models.NewStat(deviceId, 0, 0, 0)
	}
	switch action {
	case Win:
		stat.Wins++
	case Loss:
		stat.Loses++
	case Draw:
		stat.Draws++
	}

	saveStat, err := s.dao.SaveStat(*stat)
	if err != nil {
		return nil, err
	}

	return &saveStat, nil
}
