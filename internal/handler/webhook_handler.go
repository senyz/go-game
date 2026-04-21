package handler

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/senyz/go-game/interfaces"
)

type WebhookHandler struct {
	gameService interfaces.GameService
}

func NewWebhookHandler(gameService interfaces.GameService) *WebhookHandler {
	return &WebhookHandler{gameService: gameService}
}

func (h *WebhookHandler) HandleMessage(c *gin.Context) {
	var request struct {
		UserID  string `json:"user_id"`
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	// Парсим ID пользователя
	userID, err := parseUserID(request.UserID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Получаем текущий прогресс пользователя, чтобы определить текущую сцену
	currentSceneID, err := h.gameService.GetCurrentSceneID(userID)
	if err != nil {
		log.Printf("Failed to get current scene for user %d: %v", userID, err)
		c.JSON(500, gin.H{"error": "Failed to retrieve game progress"})
		return
	}

	// Обрабатываем ответ пользователя
	nextScene, err := h.gameService.ProcessAnswer(userID, currentSceneID, request.Message)
	if err != nil {
		log.Printf("Game processing error for user %d: %v", userID, err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Формируем ответ в зависимости от типа сцены
	response := gin.H{
		"scene_id": nextScene.ID,
		"title":    nextScene.Title,
		"text":     nextScene.Description,
	}

	if nextScene.Question != "" {
		response["question"] = nextScene.Question
	} else {
		// Это финальная сцена — добавляем сообщение о завершении
		response["message"] = "Поздравляем! Вы завершили приключение!"
	}

	c.JSON(200, response)
}

// parseUserID преобразует строку в uint, возвращает ошибку при неудаче
func parseUserID(idStr string) (uint, error) {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
