package repository

import (
	"context"
	"log"
	"time"

	"github.com/senyz/go-game/interfaces"
	"github.com/senyz/go-game/internal/models"
	"gorm.io/gorm"
)

type SceneRepository struct {
	db *gorm.DB
}

func NewSceneRepository(db *gorm.DB) interfaces.SceneRepository {
	return &SceneRepository{db: db}
}

func (s *SceneRepository) GetSceneByID(id uint) (*models.Scene, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var scene models.Scene
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&scene).Error; err != nil {
		log.Printf("Failed to find scene with ID %d: %v", id, err)
		return nil, err
	}

	return &scene, nil
}

func (s *SceneRepository) GetNextScene(currentSceneID uint, isCorrect bool) (*models.Scene, error) {
	currentScene, err := s.GetSceneByID(currentSceneID)
	if err != nil {
		return nil, err
	}

	var nextSceneID uint
	if isCorrect {
		if currentScene.NextSceneID == nil {
			// Это финальная сцена
			return currentScene, nil
		}
		nextSceneID = *currentScene.NextSceneID
	} else {
		if currentScene.FailureSceneID == nil {
			// Сцена неудачи не задана — возвращаем текущую
			return currentScene, nil
		}
		nextSceneID = *currentScene.FailureSceneID
	}

	nextScene, err := s.GetSceneByID(nextSceneID)
	if err != nil {
		log.Printf("Failed to get next scene (ID: %d): %v", nextSceneID, err)
		return nil, err
	}

	return nextScene, nil
}

func (s *SceneRepository) GetFirstSceneByStoryID(storyID uint) (*models.Scene, error) {
	var scene models.Scene
	err := s.db.Where("story_id = ?", storyID).Order("id ASC").First(&scene).Error
	if err != nil {
		return nil, err
	}
	return &scene, nil
}
