package db

import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "github.com/senyz/go-game/internal/config"
)

type Database struct {
    DB *gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
    db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return &Database{DB: db}, nil
}

func (d *Database) Migrate() error {
    return d.DB.AutoMigrate(
        &models.User{},
        &models.Story{},
        &models.Scene{},
        &models.UserProgress{},
    )
}
