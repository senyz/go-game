package handler

import (
	"github.com/senyz/go-game/internal/service"

	"github.com/gin-gonic/gin"
)

type StoryHandler struct {
	storyService *service.StoryService
}

func NewStoryHandler(storyService *service.StoryService) *StoryHandler {
	return &StoryHandler{storyService: storyService}
}

func (h *StoryHandler) GetStories(c *gin.Context) {
	stories, err := h.storyService.GetAllStories()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, stories)
}
