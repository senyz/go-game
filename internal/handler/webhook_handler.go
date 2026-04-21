package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/senyz/go-game/internal/service"
)

type WebhookHandler struct {
    gameService service.GameService
}

func NewWebhookHandler(gameService service.GameService) *WebhookHandler {
    return &WebhookHandler{gameService: gameService}
}

func (h *WebhookHandler) HandleMessage(c *gin.Context) {
    var request struct {
        UserID string `json:"user_id"`
        Message string `json:"message"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }

    userID := parseUserID(request.UserID) // вспомогательная функция
    response, err := h.gameService.ProcessAnswer(userID, currentSceneID, request.Message)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"response": response.Question})
}
