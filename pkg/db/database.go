package db

import (
	"github.com/senyz/go-game/internal/models"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func (d *Database) Migrate() error {
	return d.DB.AutoMigrate(
		&models.User{},
		&models.Story{},
		&models.Scene{},
		&models.UserProgress{},
	)
}
