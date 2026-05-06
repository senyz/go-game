package models

// IncomingMessage представляет входящее сообщение от пользователя
type IncomingMessage struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
	Text      string `json:"text"`
	ChatID    string `json:"chat_id"`
	Timestamp int64  `json:"timestamp"`
}
