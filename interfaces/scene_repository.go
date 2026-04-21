// interfaces/scene_repository.go
package interfaces

import "github.com/senyz/go-game/internal/models"

type SceneRepository interface {
	GetFirstSceneByStoryID(uint) (*models.Scene, error)
	GetSceneByID(id uint) (*models.Scene, error)
	GetNextScene(currentSceneID uint, isCorrect bool) (*models.Scene, error)
}
