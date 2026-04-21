package repository

import (
	"context"
	"log"
	"time"

	"github.com/senyz/go-game/interfaces"
	"github.com/senyz/go-game/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) CreateUser(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := &models.User{
		Username: username,
	}

	if err := u.db.WithContext(ctx).Create(user).Error; err != nil {
		log.Printf("failed to create user %s: %v", username, err)
		return nil, err
	}
	return user, nil
}
func (u *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := u.db.WithContext(ctx).First(user, id).Error; err != nil {
		log.Printf("Failed to find user with ID %d: %v", id, err)
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) UpdateUserProgress(userID, sceneID uint, completed bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := u.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("progress", sceneID)

	if result.Error != nil {
		log.Printf("Failed to update progress for user %d: %v", userID, result.Error)
		return result.Error
	}

	// Дополнительно обновляем таблицу прогресса пользователя
	var progress models.UserProgress
	err := u.db.Where("user_id = ? AND scene_id = ?", userID, sceneID).First(&progress).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if progress.ID == 0 {
		// Создаём новую запись прогресса
		progress = models.UserProgress{
			UserID:    userID,
			SceneID:   sceneID,
			Completed: completed,
		}
		if createErr := u.db.Create(&progress).Error; createErr != nil {
			log.Printf("Failed to create user progress: %v", createErr)
			return createErr
		}
	} else {
		// Обновляем существующую запись
		if updateErr := u.db.Model(&progress).Update("completed", completed).Error; updateErr != nil {
			log.Printf("Failed to update user progress: %v", updateErr)
			return updateErr
		}
	}

	return nil
}
