// pkg/messenger/mock_client.go
package messenger

import (
	"context"
	"sync"

	"github.com/senyz/go-game/interfaces"
)

// MockClient - мок для тестирования
type MockClient struct {
	mu            sync.RWMutex
	sentMessages  []SentMessage
	incomingQueue []*interfaces.IncomingMessage
	webhookURL    string
	shouldFail    bool
	failError     error
}

type SentMessage struct {
	UserID string
	Text   string
}

func NewMockClient() *MockClient {
	return &MockClient{
		sentMessages:  make([]SentMessage, 0),
		incomingQueue: make([]*interfaces.IncomingMessage, 0),
	}
}

// SetFailMode включает режим ошибки для тестирования
func (m *MockClient) SetFailMode(shouldFail bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = shouldFail
	m.failError = err
}

// AddIncomingMessage добавляет входящее сообщение в очередь
func (m *MockClient) AddIncomingMessage(msg *interfaces.IncomingMessage) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.incomingQueue = append(m.incomingQueue, msg)
}

// GetSentMessages возвращает все отправленные сообщения
func (m *MockClient) GetSentMessages() []SentMessage {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]SentMessage{}, m.sentMessages...)
}

// Clear очищает историю
func (m *MockClient) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sentMessages = make([]SentMessage, 0)
	m.incomingQueue = make([]*interfaces.IncomingMessage, 0)
	m.shouldFail = false
	m.failError = nil
}

// SendMessage реализация интерфейса
func (m *MockClient) SendMessage(ctx context.Context, userID, text string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldFail {
		return m.failError
	}

	m.sentMessages = append(m.sentMessages, SentMessage{
		UserID: userID,
		Text:   text,
	})
	return nil
}

// ReceiveMessage реализация интерфейса
func (m *MockClient) ReceiveMessage(ctx context.Context) (*interfaces.IncomingMessage, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldFail {
		return nil, m.failError
	}

	if len(m.incomingQueue) == 0 {
		return nil, nil // Нет сообщений
	}

	msg := m.incomingQueue[0]
	m.incomingQueue = m.incomingQueue[1:]
	return msg, nil
}

// SetWebhook реализация интерфейса
func (m *MockClient) SetWebhook(ctx context.Context, webhookURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldFail {
		return m.failError
	}

	m.webhookURL = webhookURL
	return nil
}

// GetWebhookURL возвращает установленный webhook (для тестов)
func (m *MockClient) GetWebhookURL() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.webhookURL
}
