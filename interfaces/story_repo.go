package interfaces

import "github.com/senyz/go-game/internal/models"

// StoryRepository определяет контракт для работы с историями в БД
type StoryRepository interface {
	// FindAll возвращает все доступные истории
	FindAll() ([]*models.Story, error)

	// FindByID возвращает историю по ID
	FindByID(id int) (*models.Story, error)
	// FindActiveStories возвращает только активные (доступные для игры) истории
	FindActiveStories() ([]*models.Story, error)
	// Create создаёт новую историю
	Create(story *models.Story) error
	// Update обновляет существующую историю
	Update(story *models.Story) error
	// Delete помечает историю как удалённую (soft delete) или удаляет физически
	Delete(id int) error

	// FindStoriesWithScenes возвращает истории с предварительно загруженными сценами
	FindStoriesWithScenes() ([]*models.Story, error)
	// Count возвращает общее количество историй
	Count() (int64, error)
}
