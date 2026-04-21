package logger

import (
    "github.com/sirupsen/logrus"
    "github.com/senyz/go-game/internal/config"
)

func NewLogger(cfg *config.Config) *logrus.Logger {
    log := logrus.New()

    // Формат вывода
    if cfg.Logging.Format == "json" {
        log.SetFormatter(&logrus.JSONFormatter{})
    } else {
        log.SetFormatter(&logrus.TextFormatter{})
    }

    // Уровень логирования
    level, err := logrus.ParseLevel(cfg.Logging.Level)
    if err != nil {
        level = logrus.InfoLevel
    }
    log.SetLevel(level)

    return log
}
