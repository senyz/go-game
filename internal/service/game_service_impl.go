package service

import (
	"github.com/senyz/go-game/interfaces"
	"github.com/senyz/go-game/internal/models"
)

type gameServiceImpl struct {
	userRepo  interfaces.UserRepository
	sceneRepo interfaces.SceneRepository
}

func NewGameService(userRepo interfaces.UserRepository, sceneRepo interfaces.SceneRepository) interfaces.GameService {
	return &gameServiceImpl{
		userRepo:  userRepo,
		sceneRepo: sceneRepo,
	}
}

func (s *gameServiceImpl) StartGame(userID uint, storyID uint) (*models.Scene, error) {
	// Находим первую сцену сюжета (с минимальным ID для данного story_id)
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

	var nextScene *models.Scene
	if isCorrect {
		// Переход к следующей сцене
		if currentScene.NextSceneID == nil {
			// Финальная сцена — игра завершена
			return currentScene, nil
		}
		nextScene, err = s.sceneRepo.GetSceneByID(*currentScene.NextSceneID)
		if err != nil {
			return nil, err
		}
	} else {
		// Переход к сцене неудачи
		if currentScene.FailureSceneID == nil {
			// Сцена неудачи не задана — остаёмся на текущей
			return currentScene, nil
		}
		nextScene, err = s.sceneRepo.GetSceneByID(*currentScene.FailureSceneID)
		if err != nil {
			return nil, err
		}
	}

	// Обновляем прогресс пользователя
	completed := isCorrect && nextScene.Question == "" // Завершено, если правильный ответ и нет вопроса (финальная сцена)
	if err := s.userRepo.UpdateUserProgress(userID, nextScene.ID, completed); err != nil {
		return nil, err
	}

	return nextScene, nil
}

func (s *gameServiceImpl) validateAnswer(correct, userAnswer string) bool {
	// Простая проверка ответа (можно расширить с учётом регистра, пробелов и т. д.)
	return correct == userAnswer
}

func (s *gameServiceImpl) GetHint(sceneID uint) (string, error) {
	scene, err := s.sceneRepo.GetSceneByID(sceneID)
	if err != nil {
		return "", err
	}
	return scene.Hint, nil
}

func (s *gameServiceImpl) GetCurrentSceneID(userID uint) (uint, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return 0, err
	}
	return user.Progress, nil
}
