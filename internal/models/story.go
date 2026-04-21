// models/story.go
package models

type Story struct {
    ID          uint   `gorm:"primaryKey"`
    Title       string
    Description string
    IsActive    bool
    Scenes      []Scene
}