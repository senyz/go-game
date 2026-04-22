// internal/service/messenger_service.go
package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/senyz/go-game/interfaces"
	"github.com/senyz/go-game/internal/models"
)

type MessengerService struct {
	client      interfaces.MessengerClient
	gameService interfaces.GameService
	userRepo    interfaces.UserRepository
}

func NewMessengerService(
	client interfaces.MessengerClient,
	gameService interfaces.GameService,
	userRepo interfaces.UserRepository,
) *MessengerService {
	return &MessengerService{
		client:      client,
		gameService: gameService,
		userRepo:    userRepo,
	}
}

// ProcessIncomingMessage обрабатывает входящее сообщение от пользователя
func (m *MessengerService) ProcessIncomingMessage(ctx context.Context, msg *interfaces.IncomingMessage) error {
	// 1. Находим или создаём пользователя
	userID, err := m.getOrCreateUser(msg.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Получаем текущую сцену пользователя
	currentSceneID, err := m.gameService.GetCurrentSceneID(userID)
	if err != nil {
		// Если ошибка — возможно, пользователь новый, начинаем игру
		if err.Error() == "record not found" {
			return m.startNewGame(ctx, userID, msg)
		}
		return fmt.Errorf("failed to get current scene: %w", err)
	}

	// 3. Обрабатываем ответ
	nextScene, err := m.gameService.ProcessAnswer(userID, currentSceneID, msg.Text)
	if err != nil {
		return fmt.Errorf("failed to process answer: %w", err)
	}

	// 4. Формируем ответное сообщение
	responseText := m.formatSceneMessage(nextScene)

	// 5. Отправляем ответ пользователю
	if err := m.client.SendMessage(ctx, msg.UserID, responseText); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// startNewGame начинает новую игру для пользователя
func (m *MessengerService) startNewGame(ctx context.Context, userID uint, msg *interfaces.IncomingMessage) error {
	// Начинаем игру со story_id = 1
	firstScene, err := m.gameService.StartGame(userID, 1)
	if err != nil {
		return fmt.Errorf("failed to start game: %w", err)
	}

	responseText := m.formatSceneMessage(firstScene)
	return m.client.SendMessage(ctx, msg.UserID, responseText)
}

// getOrCreateUser получает или создаёт пользователя
func (m *MessengerService) getOrCreateUser(username string) (uint, error) {
	// Пробуем найти пользователя
	users, err := m.userRepo.FindByUsername(username) // нужно добавить这个方法
	if err == nil && users != nil {
		return users.ID, nil
	}

	// Создаём нового
	user, err := m.userRepo.CreateUser(username)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

// formatSceneMessage форматирует сцену в текст для отправки
func (m *MessengerService) formatSceneMessage(scene *models.Scene) string {
	message := fmt.Sprintf("*%s*\n\n%s", scene.Title, scene.Description)

	if scene.Question != "" {
		message += fmt.Sprintf("\n\n❓ *Вопрос:* %s", scene.Question)
	}

	if scene.Hint != "" {
		message += fmt.Sprintf("\n\n💡 *Подсказка:* %s", scene.Hint)
	}

	return message
}

// RunLongPolling запускает long polling для получения сообщений
func (m *MessengerService) RunLongPolling(ctx context.Context) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			msg, err := m.client.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				continue
			}

			if msg != nil {
				log.Printf("Received message from %s: %s", msg.UserID, msg.Text)
				if err := m.ProcessIncomingMessage(ctx, msg); err != nil {
					log.Printf("Error processing message: %v", err)
					// Отправляем пользователю сообщение об ошибке
					m.client.SendMessage(ctx, msg.UserID, "❌ Ошибка! Попробуйте позже.")
				}
			}
		}
	}
}
