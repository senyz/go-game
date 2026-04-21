// repository/user_repository.go
package interfaces

type UserRepository interface {
    CreateUser(username string) (*User, error)
    GetUserByID(id uint) (*User, error)
    UpdateUserProgress(userID, sceneID uint, completed bool) error
}