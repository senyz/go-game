// pkg/messenger/max_client_test.go
package messenger

import (
	"context"
	"testing"
	"time"

	"github.com/senyz/go-game/interfaces"
)

func TestMockClient_SendMessage(t *testing.T) {
	ctx := context.Background()
	mock := NewMockClient()

	// Тест успешной отправки
	err := mock.SendMessage(ctx, "user123", "Hello, Morty!")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	messages := mock.GetSentMessages()
	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}

	if messages[0].UserID != "user123" {
		t.Errorf("Expected userID user123, got %s", messages[0].UserID)
	}

	if messages[0].Text != "Hello, Morty!" {
		t.Errorf("Expected text 'Hello, Morty!', got '%s'", messages[0].Text)
	}
}

func TestMockClient_SendMessage_Fail(t *testing.T) {
	ctx := context.Background()
	mock := NewMockClient()

	// Настраиваем режим ошибки
	expectedErr := context.DeadlineExceeded
	mock.SetFailMode(true, expectedErr)

	err := mock.SendMessage(ctx, "user123", "Hello")
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestMockClient_ReceiveMessage(t *testing.T) {
	ctx := context.Background()
	mock := NewMockClient()

	// Добавляем тестовое сообщение
	expectedMsg := &interfaces.IncomingMessage{
		UserID:    "user456",
		MessageID: "msg001",
		Text:      "2 + 2 = ?",
		ChatID:    "chat001",
		Timestamp: time.Now().Unix(),
	}
	mock.AddIncomingMessage(expectedMsg)

	// Получаем сообщение
	msg, err := mock.ReceiveMessage(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if msg == nil {
		t.Fatal("Expected message, got nil")
	}

	if msg.UserID != expectedMsg.UserID {
		t.Errorf("Expected userID %s, got %s", expectedMsg.UserID, msg.UserID)
	}

	if msg.Text != expectedMsg.Text {
		t.Errorf("Expected text %s, got %s", expectedMsg.Text, msg.Text)
	}
}

func TestMockClient_ReceiveMessage_Empty(t *testing.T) {
	ctx := context.Background()
	mock := NewMockClient()

	// Нет сообщений в очереди
	msg, err := mock.ReceiveMessage(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if msg != nil {
		t.Errorf("Expected nil message, got %v", msg)
	}
}

func TestMockClient_SetWebhook(t *testing.T) {
	ctx := context.Background()
	mock := NewMockClient()

	webhookURL := "https://example.com/webhook"
	err := mock.SetWebhook(ctx, webhookURL)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mock.GetWebhookURL() != webhookURL {
		t.Errorf("Expected webhook URL %s, got %s", webhookURL, mock.GetWebhookURL())
	}
}

func TestMockClient_Clear(t *testing.T) {
	ctx := context.Background()
	mock := NewMockClient()

	// Отправляем сообщение
	mock.SendMessage(ctx, "user1", "test")
	messages := mock.GetSentMessages()
	if len(messages) != 1 {
		t.Error("Message not sent")
	}

	// Очищаем
	mock.Clear()

	messages = mock.GetSentMessages()
	if len(messages) != 0 {
		t.Errorf("Expected 0 messages after clear, got %d", len(messages))
	}
}

// Интеграционный тест с несколькими операциями
func TestMockClient_MultipleOperations(t *testing.T) {
	ctx := context.Background()
	mock := NewMockClient()

	// Отправляем несколько сообщений
	messages := []struct {
		userID string
		text   string
	}{
		{"user1", "Message 1"},
		{"user2", "Message 2"},
		{"user3", "Message 3"},
	}

	for _, msg := range messages {
		if err := mock.SendMessage(ctx, msg.userID, msg.text); err != nil {
			t.Errorf("Failed to send message: %v", err)
		}
	}

	// Проверяем отправленные сообщения
	sent := mock.GetSentMessages()
	if len(sent) != len(messages) {
		t.Errorf("Expected %d messages, got %d", len(messages), len(sent))
	}

	// Добавляем входящее сообщение
	incoming := &interfaces.IncomingMessage{
		UserID: "user1",
		Text:   "Answer",
	}
	mock.AddIncomingMessage(incoming)

	// Получаем его
	received, err := mock.ReceiveMessage(ctx)
	if err != nil {
		t.Errorf("Receive error: %v", err)
	}

	if received.Text != "Answer" {
		t.Errorf("Expected 'Answer', got '%s'", received.Text)
	}
}
