// interfaces/messenger_client.go
package interfaces

import (
	"context"

	"github.com/senyz/go-game/internal/models"
)

// MessengerClient определяет контракт для работы с любым мессенджером
type MessengerClient interface {
	// SendMessage отправляет сообщение пользователю
	SendMessage(ctx context.Context, userID string, text string) error

	// ReceiveMessage получает сообщение от пользователя (для Long Polling)
	ReceiveMessage(ctx context.Context) (*models.IncomingMessage, error)

	// SetWebhook устанавливает webhook URL для получения сообщений (для production)
	SetWebhook(ctx context.Context, webhookURL string) error
}
