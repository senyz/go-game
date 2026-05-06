// interfaces/game_service.go
package interfaces

import "github.com/senyz/go-game/internal/models"

type GameService interface {
	StartGame(userID uint, storyID uint) (*models.Scene, error)
	ProcessAnswer(userID uint, sceneID uint, answer string) (*models.Scene, error)
	GetHint(sceneID uint) (string, error)
	GetCurrentSceneID(userID uint) (uint, error)
	GetUserProgress(userID uint) ([]models.UserProgress, error)
	CheckAnswer(userID uint, sceneID uint, answer string) (isCorrect bool, nextSceneID uint, err error)
}
