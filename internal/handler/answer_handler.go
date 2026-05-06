package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *WebhookHandler) SubmitAnswer(c *gin.Context) {
	var req struct {
		UserID  uint   `json:"user_id"`
		SceneID uint   `json:"scene_id"`
		Answer  string `json:"answer"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isCorrect, nextScene, err := h.gameService.CheckAnswer(req.UserID, req.SceneID, req.Answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"correct":    isCorrect,
		"next_scene": nextScene,
	})
}
