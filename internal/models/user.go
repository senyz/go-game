// models/user.go
package models

type User struct {
    ID        uint   `gorm:"primaryKey"`
    Username  string `gorm:"unique"`
    Progress  uint   // ID текущей сцены
    Score     int
    CreatedAt time.Time
}
