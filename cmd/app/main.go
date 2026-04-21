package main

import (
	"fmt"
	"log"
	"time"

	"github.com/senyz/go-game/interfaces"
	"github.com/senyz/go-game/internal/config"
	"github.com/senyz/go-game/internal/migration"
	"github.com/senyz/go-game/internal/repository"
	"github.com/senyz/go-game/pkg/db"
	"github.com/senyz/go-game/pkg/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализируем логирование
	appLogger := logger.NewLogger(cfg)
	appLogger.Info("Logger initialized")

	// Подключаемся к БД
	database, err := db.NewDatabase(cfg)
	if err != nil {
		appLogger.Fatalf("Failed to connect to database: %v", err)
	}
	appLogger.Info("Database connected")

	// Выполняем миграции
	migrator := migration.NewMigrator(database.DB)
	if err := migrator.Run(); err != nil {
		appLogger.Fatalf("Migrations failed: %v", err)
	}
	appLogger.Info("All migrations applied successfully")

	// Создаём репозитории после успешного применения миграций
	userRepo := repository.NewUserRepository(database.DB)
	sceneRepo := repository.NewSceneRepository(database.DB)

	// Тестирование репозиториев (опционально, для отладки)
	testRepositories(appLogger, userRepo, sceneRepo)

	fmt.Println("Application is ready to start!")

	// Здесь будет инициализация Gin, сервисов, обработчиков и запуск сервера
	// startServer(appLogger, userRepo, sceneRepo) // пример вызова запуска сервера
}

// testRepositories — вспомогательная функция для тестирования репозиториев
func testRepositories(logger *logrus.Logger, userRepo interfaces.UserRepository, sceneRepo interfaces.SceneRepository) {
	logger.Info("Starting repository tests...")

	// Тест создания пользователя
	user, err := userRepo.CreateUser("test_user_" + fmt.Sprintf("%d", time.Now().Unix()))
	if err != nil {
		logger.Errorf("Failed to create test user: %v", err)
	} else {
		logger.Infof("Created test user with ID: %d", user.ID)
	}

	// Тест получения сцены (с проверкой существования)
	scene, err := sceneRepo.GetSceneByID(1)
	if err != nil {
		logger.Warnf("Failed to get scene with ID 1: %v (this is OK if no test data)", err)
	} else {
		logger.Infof("Got scene: %s, question: %s", scene.Title, scene.Question)

		// Тест перехода к следующей сцене (если сцена найдена)
		nextScene, err := sceneRepo.GetNextScene(1, true)
		if err != nil {
			logger.Errorf("Failed to get next scene: %v", err)
		} else {
			logger.Infof("Next scene: %s", nextScene.Title)
		}
	}
}
