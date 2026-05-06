package repository

import (
	"github.com/senyz/go-game/interfaces"
	"github.com/senyz/go-game/internal/models"
	"gorm.io/gorm"
)

// StoryRepository определяет контракт для работы с историями в БД
type StoryRepository struct {
	db *gorm.DB
}

func NewStoryRepository(db *gorm.DB) interfaces.StoryRepository {
	return &StoryRepository{db: db}
}

// FindAll возвращает все доступные истории
func (S *StoryRepository) FindAll() ([]*models.Story, error) {
	stories := make([]*models.Story, 0)
	if err := S.db.Find(&stories).Error; err != nil {
		return nil, err
	}
	return stories, nil
}

// FindByID возвращает историю по ID
func (S *StoryRepository) FindByID(id int) (*models.Story, error) {
	if err := S.db.First(&models.Story{}, id).Error; err != nil {
		return nil, err
	}
	return &models.Story{}, nil

}

// FindActiveStories возвращает только активные (доступные для игры) истории
func (S *StoryRepository) FindActiveStories() ([]*models.Story, error) {
	stories := make([]*models.Story, 0)
	if err := S.db.Where("is_active = ?", true).Find(&stories).Error; err != nil {
		return nil, err
	}
	return stories, nil

}

// Create создаёт новую историю
func (S *StoryRepository) Create(story *models.Story) error {
	if err := S.db.Create(story).Error; err != nil {
		return err
	}

	return nil
}

// Update обновляет существующую историю
func (S *StoryRepository) Update(story *models.Story) error {
	if err := S.db.Save(story).Error; err != nil {
		return err
	}
	return nil
}

// Delete помечает историю как удалённую (soft delete) или удаляет физически
func (S *StoryRepository) Delete(id int) error {
	if err := S.db.Delete(&models.Story{}, id).Error; err != nil {
		return err
	}
	return nil
}

// FindStoriesWithScenes возвращает истории с предварительно загруженными сценами
func (S *StoryRepository) FindStoriesWithScenes() ([]*models.Story, error) {
	stories := make([]*models.Story, 0)
	if err := S.db.Preload("Scenes").Find(&stories).Error; err != nil {
		return nil, err
	}
	return stories, nil
}

// Count возвращает общее количество историй
func (S *StoryRepository) Count() (int64, error) {
	var count int64
	if err := S.db.Model(&models.Story{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
