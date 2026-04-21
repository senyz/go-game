// models/user_progress.go
package models

type UserProgress struct {
    ID       uint   `gorm:"primaryKey"`
    UserID   uint
    SceneID  uint
    Attempts int
    Completed bool
}