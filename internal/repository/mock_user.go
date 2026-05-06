package repository

import (
	"fmt"

	"github.com/senyz/go-game/internal/models"
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

func (m *MockUserRepository) GetUserProgress(userID uint) ([]models.UserProgress, error) {
	user, ok := m.users[userID]
	if !ok {
		return nil, fmt.Errorf("user not found: %d", userID)
	}

	// Если прогресс пользователя не установлен (Progress == 0), возвращаем пустой срез
	if user.Progress == 0 {
		return nil, nil
	}

	// Формируем прогресс пользователя: текущая сцена с флагом IsCompleted = false
	progress := []models.UserProgress{
		{
			SceneID:     user.Progress,
			IsCompleted: false,
		},
	}

	return progress, nil
}

func (m *MockUserRepository) UpdateUserProgress(userID, sceneID uint, completed bool) error {
	if user, ok := m.users[userID]; ok {
		user.Progress = sceneID
		return nil
	}
	return nil
}

func (m *MockUserRepository) GetUserByUsername(username string) (*models.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}
