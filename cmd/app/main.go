package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/senyz/go-game/interfaces"
	"github.com/senyz/go-game/internal/config"
	"github.com/senyz/go-game/internal/handler"
	"github.com/senyz/go-game/internal/migration"
	"github.com/senyz/go-game/internal/repository"
	"github.com/senyz/go-game/internal/service"
	"github.com/senyz/go-game/pkg/db"
	"github.com/senyz/go-game/pkg/logger"
	"github.com/senyz/go-game/pkg/messenger"
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

	// Выполняем миграции с указанием диалекта
	migrator := migration.NewMigrator(database, "./migrations", cfg.Database.Dialect)
	if err := migrator.Run(); err != nil {
		appLogger.Fatalf("Migrations failed: %v", err)
	}
	appLogger.Info("All migrations applied successfully")

	// Создаём репозитории после успешного применения миграций
	userRepo := repository.NewUserRepository(database)
	sceneRepo := repository.NewSceneRepository(database)

	// Тестирование репозиториев (опционально, для отладки)
	testRepositories(appLogger, userRepo, sceneRepo)

	fmt.Println("Application is ready to start!")

	// Передаём cfg в startServer
	startServer(appLogger, userRepo, sceneRepo, cfg) // ← добавили cfg
}

// Обновляем сигнатуру функции startServer
func startServer(appLogger *logrus.Logger, userRepo interfaces.UserRepository, sceneRepo interfaces.SceneRepository, cfg *config.Config) {
	// 1. Создаём сервис игры
	gameService := service.NewGameService(userRepo, sceneRepo)

	// 2. Создаём клиент мессенджера
	var messengerClient interfaces.MessengerClient
	if cfg.Messenger.UseWebhook {
		messengerClient = messenger.NewMAXClient(cfg.Messenger.APIURL, cfg.Messenger.Token)
		if err := messengerClient.SetWebhook(context.Background(), cfg.Messenger.WebhookURL); err != nil {
			appLogger.Warnf("Failed to set webhook: %v", err)
		}
	} else {
		messengerClient = messenger.NewMAXClient(cfg.Messenger.APIURL, cfg.Messenger.Token)
	}

	// 3. Создаём сервис мессенджера
	messengerService := service.NewMessengerService(messengerClient, gameService, userRepo)

	// 4. Создаём HTTP обработчик
	webhookHandler := handler.NewWebhookHandler(gameService)

	// 5. Настраиваем Gin
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.POST("/webhook", webhookHandler.HandleMessage)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 6. Создаём HTTP сервер
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: router,
	}

	// 7. Запускаем long polling (если не используется webhook)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if !cfg.Messenger.UseWebhook {
		go func() {
			appLogger.Info("Starting long polling for messages...")
			if err := messengerService.RunLongPolling(ctx); err != nil {
				appLogger.Errorf("Long polling error: %v", err)
			}
		}()
	}

	// 8. Graceful shutdown
	go func() {
		appLogger.Infof("Server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Отменяем контекст для long polling
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		appLogger.Fatalf("Server forced to shutdown: %v", err)
	}

	appLogger.Info("Server exited gracefully")
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
