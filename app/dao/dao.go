package dao

import (
	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Dao interface {
	GameDao
}

type dao struct {
	db *gorm.DB
}

func NewDao() (Dao, error) {
	conn, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "a serious error occurred connecting to the database")
	}

	daoI := &dao{
		db: conn,
	}

	err = daoI.migrateAll()
	if err != nil {
		return nil, errors.Wrap(err, "a serious error occurred migrating the databases")
	}

	return daoI, nil
}

func (d dao) migrateAll() error {
	toMigrate := []any{
		&models.Game{},
		&models.Stat{},
	}

	for _, model := range toMigrate {
		err := d.db.AutoMigrate(model)
		if err != nil {
			return errors.Wrap(err, "a serious error occurred migrating the database")
		}
	}

	return nil
}

func (d dao) getDb(tx ...*gorm.DB) *gorm.DB {
	db := d.db
	if len(tx) == 1 {
		db = tx[0]
	}
	return db
}
