package service

type GameService interface {
    StartGame(userID uint, storyID uint) (*models.Scene, error)
    ProcessAnswer(userID uint, sceneID uint, answer string) (*models.Scene, error)
    GetHint(sceneID uint) (string, error)
}


type gameServiceImpl struct {
    userRepo repository.UserRepository
    sceneRepo repository.SceneRepository
}

func NewGameService(userRepo repository.UserRepository, sceneRepo repository.SceneRepository) GameService {
    return &gameServiceImpl{
        userRepo: userRepo,
        sceneRepo: sceneRepo,
    }
}

func (s *gameServiceImpl) StartGame(userID uint, storyID uint) (*models.Scene, error) {
    // Логика старта игры: находим первую сцену сюжета
    firstScene, err := s.sceneRepo.GetFirstSceneByStoryID(storyID)
    if err != nil {
        return nil, err
    }

    // Обновляем прогресс пользователя
    if err := s.userRepo.UpdateUserProgress(userID, firstScene.ID, false); err != nil {
        return nil, err
    }

    return firstScene, nil
}

func (s *gameServiceImpl) ProcessAnswer(userID uint, sceneID uint, answer string) (*models.Scene, error) {
    currentScene, err := s.sceneRepo.GetSceneByID(sceneID)
    if err != nil {
        return nil, err
    }

    isCorrect := s.validateAnswer(currentScene.CorrectAnswer, answer)

    if isCorrect {
        // Переход к следующей сцене
        nextScene, err := s.sceneRepo.GetNextScene(currentScene.NextSceneID)
        if err != nil {
            return nil, err
        }
        s.userRepo.UpdateUserProgress(userID, nextScene.ID, true)
        return nextScene, nil
    } else {
        // Переход к сцене неудачи
        failureScene, err := s.sceneRepo.GetSceneByID(*currentScene.FailureSceneID)
        if err != nil {
            return nil, err
        }
        return failureScene, nil
    }
}

func (s *gameServiceImpl) validateAnswer(correct, userAnswer string) bool {
    // Простая проверка ответа (можно расширить с учётом регистра, пробелов и т. д.)
    return correct == userAnswer
}