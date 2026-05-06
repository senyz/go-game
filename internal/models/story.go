// models/story.go
package models

type Story struct {
	ID          int    `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string
	IsActive    bool    `gorm:"default:true"`
	Scenes      []Scene `gorm:"foreignKey:StoryID"`
}

type Option struct {
	ID          int    `gorm:"primaryKey"`
	SceneID     int    `gorm:"not null"`
	Text        string `gorm:"not null"`
	NextSceneID int
}
