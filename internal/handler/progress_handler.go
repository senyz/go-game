package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *WebhookHandler) GetProgress(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	progress, err := h.gameService.GetUserProgress(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}
