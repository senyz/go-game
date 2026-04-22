// interfaces/messenger_client.go
package interfaces

import "context"

// MessengerClient определяет контракт для работы с любым мессенджером
type MessengerClient interface {
	// SendMessage отправляет сообщение пользователю
	SendMessage(ctx context.Context, userID string, text string) error

	// ReceiveMessage получает сообщение от пользователя (для Long Polling)
	ReceiveMessage(ctx context.Context) (*IncomingMessage, error)

	// SetWebhook устанавливает webhook URL для получения сообщений (для production)
	SetWebhook(ctx context.Context, webhookURL string) error
}

// IncomingMessage представляет входящее сообщение от пользователя
type IncomingMessage struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
	Text      string `json:"text"`
	ChatID    string `json:"chat_id"`
	Timestamp int64  `json:"timestamp"`
}
