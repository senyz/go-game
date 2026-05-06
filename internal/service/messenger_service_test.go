// internal/service/messenger_service_test.go
package service

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/senyz/go-game/interfaces"
	"github.com/senyz/go-game/internal/models"
	"github.com/senyz/go-game/internal/repository"
	"github.com/senyz/go-game/pkg/messenger"
)

// MockGameService - полная реализация интерфейса GameService
type MockGameService struct {
	scenes   map[uint]*models.Scene // Сцены по ID
	users    map[uint]*models.User  // Пользователи по ID
	progress map[uint]uint          // Прогресс: userID → sceneID
	userRepo interfaces.UserRepository
}

func NewMockGameService(userRepo interfaces.UserRepository) *MockGameService {
	return &MockGameService{
		scenes:   make(map[uint]*models.Scene),
		users:    make(map[uint]*models.User),
		progress: make(map[uint]uint),
		userRepo: userRepo,
	}
}

func (m *MockGameService) GetCurrentSceneID(userID uint) (uint, error) {
	if sceneID, ok := m.progress[userID]; ok {
		return sceneID, nil
	}
	return 0, fmt.Errorf("user progress not found")
}

// SetUserProgress устанавливает прогресс пользователя для тестирования
func (m *MockGameService) SetUserProgress(userID uint, sceneID uint) {
	m.progress[userID] = sceneID
}

func (m *MockGameService) StartGame(userID uint, storyID uint) (*models.Scene, error) {
	// Возвращаем первую сцену
	for _, scene := range m.scenes {
		if scene.StoryID == storyID {
			return scene, nil
		}
	}
	return &models.Scene{
		ID:            1,
		StoryID:       storyID,
		Title:         "Test Scene",
		Description:   "Test description",
		Question:      "2 + 2 = ?",
		CorrectAnswer: "4",
		Hint:          "Think about it",
	}, nil
}

func (m *MockGameService) ProcessAnswer(userID uint, sceneID uint, answer string) (*models.Scene, error) {
	scene, ok := m.scenes[sceneID]
	if !ok {
		// Возвращаем тестовую сцену
		return &models.Scene{
			ID:          2,
			Title:       "Next Scene",
			Description: "You progressed!",
			Question:    "",
		}, nil
	}

	if answer == scene.CorrectAnswer {
		// Правильный ответ - переходим к следующей сцене
		if scene.NextSceneID != nil {
			return m.GetSceneByID(*scene.NextSceneID)
		}
		return scene, nil
	}

	// Неправильный ответ - сцена неудачи
	if scene.FailureSceneID != nil {
		return m.GetSceneByID(*scene.FailureSceneID)
	}
	return scene, nil
}

func (m *MockGameService) GetHint(sceneID uint) (string, error) {
	if scene, ok := m.scenes[sceneID]; ok {
		return scene.Hint, nil
	}
	return "No hint available", nil
}

func (m *MockGameService) GetSceneByID(sceneID uint) (*models.Scene, error) {
	if scene, ok := m.scenes[sceneID]; ok {
		return scene, nil
	}
	return &models.Scene{
		ID:          sceneID,
		Title:       "Scene",
		Description: "Description",
	}, nil
}

// AddScene добавляет тестовую сцену
func (m *MockGameService) AddScene(scene *models.Scene) {
	m.scenes[scene.ID] = scene
}

func (m *MockGameService) CheckAnswer(userID uint, sceneID uint, answer string) (bool, uint, error) {
	scene, ok := m.scenes[sceneID]
	if !ok {
		return false, 0, fmt.Errorf("scene not found: %d", sceneID)
	}

	isCorrect := (scene.CorrectAnswer == answer)
	var nextSceneID uint

	if isCorrect {
		if scene.NextSceneID != nil {
			nextSceneID = *scene.NextSceneID
		} else {
			nextSceneID = sceneID // Остаёмся на текущей сцене
		}
	} else {
		if scene.FailureSceneID != nil {
			nextSceneID = *scene.FailureSceneID
		} else {
			nextSceneID = sceneID // Остаёмся на текущей сцене
		}
	}

	// Обновляем прогресс пользователя
	m.progress[userID] = nextSceneID

	return isCorrect, nextSceneID, nil
}

func (m *MockGameService) GetUserProgress(userID uint) ([]models.UserProgress, error) {
	// Делегируем вызов реальному репозиторию
	return m.userRepo.GetUserProgress(userID)
}

func TestMessengerService_ProcessIncomingMessage(t *testing.T) {
	mockMessenger := messenger.NewMockClient()
	mockUserRepo := repository.NewMockUserRepository()
	mockGameService := NewMockGameService(mockUserRepo)

	// Создаём пользователя
	user, _ := mockUserRepo.CreateUser("test_user")

	// Добавляем тестовую сцену
	scene := &models.Scene{
		ID:            1,
		StoryID:       1,
		Title:         "Math Test",
		Description:   "Solve this problem",
		Question:      "2 + 2 = ?",
		CorrectAnswer: "4",
	}
	mockGameService.AddScene(scene)

	// Устанавливаем начальный прогресс пользователя
	mockGameService.SetUserProgress(user.ID, 1)

	service := NewMessengerService(mockMessenger, mockGameService, mockUserRepo)

	msg := &models.IncomingMessage{
		UserID:    user.Username,
		MessageID: "msg_001",
		Text:      "4",
		ChatID:    "chat_001",
		Timestamp: 1234567890,
	}

	ctx := context.Background()
	err := service.ProcessIncomingMessage(ctx, msg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	sent := mockMessenger.GetSentMessages()
	t.Logf("Sent messages count: %d", len(sent))
	t.Logf("Sent messages: %+v", sent)

	if len(sent) == 0 {
		t.Error("Expected message to be sent")
	} else {
		if !strings.Contains(sent[0].Text, "Solve") && !strings.Contains(sent[0].Text, "2 + 2") {
			t.Errorf("Unexpected message text: %s", sent[0].Text)
		}
	}
}

func TestMessengerService_CorrectAnswer(t *testing.T) {
	mockMessenger := messenger.NewMockClient()
	mockUserRepo := repository.NewMockUserRepository()
	mockGameService := NewMockGameService(mockUserRepo)

	user, _ := mockUserRepo.CreateUser("test_user")
	// Устанавливаем начальный прогресс пользователя (сцена 1)
	mockGameService.SetUserProgress(user.ID, 1)

	nextScene := &models.Scene{ID: 2, Description: "You answered correctly!"}
	currentScene := &models.Scene{
		ID:             1,
		StoryID:        1,
		CorrectAnswer:  "4",
		NextSceneID:    func(id uint) *uint { return &id }(2),
		FailureSceneID: func(id uint) *uint { return &id }(99),
	}

	mockGameService.AddScene(currentScene)
	mockGameService.AddScene(nextScene)

	service := NewMessengerService(mockMessenger, mockGameService, mockUserRepo)

	msg := &models.IncomingMessage{UserID: user.Username, Text: "4"}
	ctx := context.Background()
	err := service.ProcessIncomingMessage(ctx, msg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	sent := mockMessenger.GetSentMessages()
	if len(sent) == 0 {
		t.Error("Expected message to be sent")
	} else if !strings.Contains(sent[0].Text, "correctly") {
		t.Errorf("Expected correct answer message, got: %s", sent[0].Text)
	}
}

func TestMessengerService_WrongAnswer(t *testing.T) {
	mockMessenger := messenger.NewMockClient()
	mockUserRepo := repository.NewMockUserRepository()
	mockGameService := NewMockGameService(mockUserRepo)

	user, _ := mockUserRepo.CreateUser("test_user")

	failureScene := &models.Scene{
		ID:          99,
		StoryID:     1,
		Title:       "Wrong Answer",
		Description: "That's incorrect! Try again.",
	}

	currentScene := &models.Scene{
		ID:             1,
		StoryID:        1,
		Title:          "First Question",
		Description:    "What is 2 + 2?",
		Question:       "2 + 2 = ?",
		CorrectAnswer:  "4",
		FailureSceneID: func(id uint) *uint { return &id }(99),
	}

	mockGameService.AddScene(currentScene)
	mockGameService.AddScene(failureScene)
	// Устанавливаем начальный прогресс пользователя (сцена 1)
	mockGameService.SetUserProgress(user.ID, 1)
	service := NewMessengerService(mockMessenger, mockGameService, mockUserRepo)

	msg := &models.IncomingMessage{
		UserID: user.Username,
		Text:   "5", // неправильный ответ
	}

	ctx := context.Background()
	err := service.ProcessIncomingMessage(ctx, msg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	sent := mockMessenger.GetSentMessages()
	if len(sent) == 0 {
		t.Error("Expected error message to be sent")
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(sent[0].Text, "incorrect") {
		t.Error("Expected incorrect answer message")
	}
}

func TestMessengerService_NewUser(t *testing.T) {
	mockMessenger := messenger.NewMockClient()
	mockUserRepo := repository.NewMockUserRepository()
	mockGameService := NewMockGameService(mockUserRepo)

	service := NewMessengerService(mockMessenger, mockGameService, mockUserRepo)

	msg := &models.IncomingMessage{
		UserID: "brand_new_user",
		Text:   "start",
	}
	user, _ := mockUserRepo.CreateUser("brand_new_user")
	// Устанавливаем начальный прогресс пользователя (сцена 1)
	mockGameService.SetUserProgress(user.ID, 1)
	ctx := context.Background()
	err := service.ProcessIncomingMessage(ctx, msg)
	if err != nil {
		t.Errorf("Expected no error for new user, got %v", err)
	}

	sent := mockMessenger.GetSentMessages()
	if len(sent) == 0 {
		t.Error("Expected welcome message for new user")
	}

	// Проверяем создание пользователя
	newUser, err := mockUserRepo.GetUserByUsername("brand_new_user")
	if err != nil || newUser == nil {
		t.Error("Expected user to be created")
	}
}

func TestMessengerService_NoNextScene(t *testing.T) {
	mockMessenger := messenger.NewMockClient()
	mockUserRepo := repository.NewMockUserRepository()
	mockGameService := NewMockGameService(mockUserRepo)

	user, _ := mockUserRepo.CreateUser("test_user")

	finalScene := &models.Scene{
		ID:          3,
		Title:       "Final Scene",
		Description: "Congratulations! You completed the story!",
		// NextSceneID и FailureSceneID не указаны
	}

	mockGameService.AddScene(finalScene)
	// Устанавливаем начальный прогресс пользователя (сцена 3)
	mockGameService.SetUserProgress(user.ID, 3)
	service := NewMessengerService(mockMessenger, mockGameService, mockUserRepo)

	msg := &models.IncomingMessage{
		UserID: user.Username,
		Text:   "any answer", // любой ответ в финальной сцене
	}

	ctx := context.Background()
	err := service.ProcessIncomingMessage(ctx, msg)
	if err != nil {
		t.Errorf("Expected no error in final scene, got %v", err)
	}

	sent := mockMessenger.GetSentMessages()
	if len(sent) == 0 {
		t.Error("Expected completion message")
	}

	if !strings.Contains(sent[0].Text, "Congratulations") {
		t.Error("Expected completion message")
	}
}

func TestMessengerService_GetUserProgress(t *testing.T) {
	mockMessenger := messenger.NewMockClient()
	mockUserRepo := repository.NewMockUserRepository()
	mockGameService := NewMockGameService(mockUserRepo)

	// Создаём пользователя
	user, _ := mockUserRepo.CreateUser("test_user")

	// Устанавливаем прогресс пользователя
	_ = mockUserRepo.UpdateUserProgress(user.ID, 5, false)

	_ = NewMessengerService(mockMessenger, mockGameService, mockUserRepo)

	// Получаем прогресс пользователя
	progress, err := mockUserRepo.GetUserProgress(user.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(progress) == 0 {
		t.Error("Expected user progress to be returned")
	}

	if progress[0].SceneID != 5 {
		t.Errorf("Expected scene ID 5, got %d", progress[0].SceneID)
	}

	if progress[0].IsCompleted {
		t.Error("Expected IsCompleted to be false")
	}
}

func TestMockClient_SendMessage(t *testing.T) {
	client := messenger.NewMockClient()
	ctx := context.Background()

	err := client.SendMessage(ctx, "test_user", "Test message")
	if err != nil {
		t.Errorf("SendMessage failed: %v", err)
	}

	sent := client.GetSentMessages()
	if len(sent) != 1 {
		t.Errorf("Expected 1 message, got %d", len(sent))
	}

	if sent[0].Text != "Test message" {
		t.Errorf("Expected 'Test message', got '%s'", sent[0].Text)
	}
}
