// models/scene.go
package models

type Scene struct {
    ID           uint   `gorm:"primaryKey"`
    StoryID      uint
    Title        string
    Description  string
    Question     string
    CorrectAnswer string
    Hint         string
    NextSceneID  *uint // nil для финальных сцен
    FailureSceneID *uint
}