// internal/service/messenger_service_test.go
package service

import (
	"context"
	"testing"

	"github.com/senyz/go-game/interfaces"
	"github.com/senyz/go-game/internal/models"
	"github.com/senyz/go-game/pkg/messenger"
)

// MockUserRepository - полная реализация интерфейса UserRepository
type MockUserRepository struct {
	users  map[uint]*models.User
	nextID uint
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[uint]*models.User),
		nextID: 1,
	}
}

func (m *MockUserRepository) CreateUser(username string) (*models.User, error) {
	user := &models.User{
		ID:       m.nextID,
		Username: username,
		Progress: 0,
		Score:    0,
	}
	m.users[user.ID] = user
	m.nextID++
	return user, nil
}

func (m *MockUserRepository) GetUserByID(id uint) (*models.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, nil // record not found
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, nil // not found
}

func (m *MockUserRepository) UpdateUserProgress(userID, sceneID uint, completed bool) error {
	if user, ok := m.users[userID]; ok {
		user.Progress = sceneID
		return nil
	}
	return nil
}

// MockGameService - полная реализация интерфейса GameService
type MockGameService struct {
	scenes         map[uint]*models.Scene
	currentSceneID uint
}

func NewMockGameService() *MockGameService {
	return &MockGameService{
		scenes: make(map[uint]*models.Scene),
	}
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

func (m *MockGameService) GetCurrentSceneID(userID uint) (uint, error) {
	return m.currentSceneID, nil
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

func TestMessengerService_ProcessIncomingMessage(t *testing.T) {
	// Создаём моки
	mockMessenger := messenger.NewMockClient()
	mockUserRepo := NewMockUserRepository()
	mockGameService := NewMockGameService()

	// Добавляем тестовую сцену
	mockGameService.AddScene(&models.Scene{
		ID:            1,
		StoryID:       1,
		Title:         "Math Test",
		Description:   "Solve this problem",
		Question:      "2 + 2 = ?",
		CorrectAnswer: "4",
		Hint:          "Add two and two",
		NextSceneID:   nil,
	})

	service := NewMessengerService(mockMessenger, mockGameService, mockUserRepo)

	// Тестовое сообщение
	msg := &interfaces.IncomingMessage{
		UserID:    "test_user_123",
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

	// Проверяем, что сообщение было отправлено
	sent := mockMessenger.GetSentMessages()
	if len(sent) == 0 {
		t.Error("Expected message to be sent")
	}

	t.Logf("Sent messages: %+v", sent)
}

func TestMessengerService_NewUser(t *testing.T) {
	mockMessenger := messenger.NewMockClient()
	mockUserRepo := NewMockUserRepository()
	mockGameService := NewMockGameService()

	service := NewMessengerService(mockMessenger, mockGameService, mockUserRepo)

	// Новый пользователь
	msg := &interfaces.IncomingMessage{
		UserID: "brand_new_user",
		Text:   "start",
	}

	ctx := context.Background()
	err := service.ProcessIncomingMessage(ctx, msg)
	if err != nil {
		t.Errorf("Expected no error for new user, got %v", err)
	}

	// Проверяем, что пользователь был создан
	sent := mockMessenger.GetSentMessages()
	if len(sent) == 0 {
		t.Error("Expected welcome message for new user")
	}
}

func TestMessengerService_WrongAnswer(t *testing.T) {
	mockMessenger := messenger.NewMockClient()
	mockUserRepo := NewMockUserRepository()
	mockGameService := NewMockGameService()

	// Создаём пользователя
	user, _ := mockUserRepo.CreateUser("test_user")

	// Добавляем сцену с неправильным ответом
	failureScene := &models.Scene{
		ID:          99,
		Title:       "Wrong Answer",
		Description: "That's incorrect! Try again.",
		Question:    "Try again: 2 + 2 = ?",
	}

	mockGameService.AddScene(&models.Scene{
		ID:             1,
		Title:          "First Question",
		Description:    "What is 2 + 2?",
		Question:       "2 + 2 = ?",
		CorrectAnswer:  "4",
		FailureSceneID: &failureScene.ID,
	})
	mockGameService.AddScene(failureScene)

	service := NewMessengerService(mockMessenger, mockGameService, mockUserRepo)

	// Отправляем неправильный ответ
	msg := &interfaces.IncomingMessage{
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

	// Проверяем, что сообщение содержит информацию об ошибке
	if len(sent) > 0 && sent[0].Text == "" {
		t.Error("Expected non-empty error message")
	}
}
