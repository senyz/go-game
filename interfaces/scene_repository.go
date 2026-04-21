// repository/scene_repository.go
package interfaces

type SceneRepository interface {
    GetSceneByID(id uint) (*Scene, error)
    GetNextScene(currentSceneID uint, isCorrect bool) (*Scene, error)
}