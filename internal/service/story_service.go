package service

import (
	"github.com/senyz/go-game/interfaces"
	"github.com/senyz/go-game/internal/models"
)

type StoryService struct {
	storyRepo interfaces.StoryRepository
}

func NewStoryService(storyRepo interfaces.StoryRepository) *StoryService {
	return &StoryService{storyRepo: storyRepo}
}

func (s *StoryService) GetAllStories() ([]*models.Story, error) {
	return s.storyRepo.FindAll()
}
