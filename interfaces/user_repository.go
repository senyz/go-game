// repository/user_repository.go
package interfaces

import "github.com/senyz/go-game/internal/models"

type UserRepository interface {
	CreateUser(username string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUserProgress(userID, sceneID uint, completed bool) error
}
