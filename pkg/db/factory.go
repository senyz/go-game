package db

import (
	"fmt"

	"github.com/senyz/go-game/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	switch cfg.Database.Dialect {
	case "postgres":
		return gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	case "mysql":
		return gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database dialect: %s", cfg.Database.Dialect)
	}
}
