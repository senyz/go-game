package models

import "time"

// Progress представляет прогресс пользователя в конкретной истории
type Progress struct {
	ID          int        `gorm:"primaryKey"`
	UserID      int        `gorm:"not null;index"`
	StoryID     int        `gorm:"not null;index"`
	SceneID     int        `gorm:"not null"`      // Текущая сцена
	IsCompleted bool       `gorm:"default:false"` // Завершена ли история
	CompletedAt *time.Time `gorm:"null"`          // Когда завершена (если завершена)
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
}
