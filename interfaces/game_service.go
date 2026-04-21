// service/game_service.go
package interfaces

type GameService interface {
    StartGame(userID uint, storyID uint) (*Scene, error)
    ProcessAnswer(userID uint, sceneID uint, answer string) (*Scene, error)
    GetHint(sceneID uint) (string, error)
}
