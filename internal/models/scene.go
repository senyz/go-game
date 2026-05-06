// models/scene.go
package models

import "time"

// type Scene struct {
// 	ID             uint `gorm:"primaryKey"`
// 	StoryID        uint
// 	Title          string
// 	Description    string
// 	Question       string
// 	CorrectAnswer  string
// 	Hint           string
// 	NextSceneID    *uint // nil для финальных сцен
// 	FailureSceneID *uint
// }

type Scene struct {
	ID             uint   `gorm:"primaryKey"`
	StoryID        uint   `gorm:"not null"`
	Title          string `gorm:"not null"`
	Description    string `gorm:"not null"`
	Question       string
	CorrectAnswer  string `gorm:"not null"`
	NextSceneID    *uint  // ID следующей сцены при правильном ответе
	FailureSceneID *uint  // ID сцены при неправильном ответе
	Hint           string
	IsFirstScene   bool      // Флаг начальной сцены истории
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}
