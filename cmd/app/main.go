package main

import (
    "fmt"
    "log"
    "github.com/senyz/go-game/internal/config"
    "github.com/senyz/go-game/pkg/db"
    "github.com/senyz/go-game/pkg/logger"
)

func main() {
    // Загружаем конфигурацию
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Инициализируем логирование
    log := logger.NewLogger(cfg)
    log.Info("Logger initialized")

    // Подключаемся к БД
    database, err := db.NewDatabase(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    log.Info("Database connected")

    // Выполняем миграции
    if err := database.Migrate(); err != nil {
        log.Fatalf("Migration failed: %v", err)
    }
    log.Info("Migrations applied")
	migrator := migration.NewMigrator(database.DB)
if err := migrator.Run(); err != nil {
    log.Fatalf("Migrations failed: %v", err)
}
log.Info("All migrations applied successfully")

    fmt.Println("Application is ready to start!")
    // Здесь будет инициализация Gin, сервисов, обработчиков и запуск сервера
}
