// pkg/messenger/max_client.go
package messenger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/senyz/go-game/interfaces"
)

// MAXClient реализует интерфейс MessengerClient для MAX API
type MAXClient struct {
	apiURL     string
	token      string
	webhookURL string
	httpClient *http.Client
}

// NewMAXClient создаёт новый экземпляр клиента для MAX
func NewMAXClient(apiURL, token string) *MAXClient {
	return &MAXClient{
		apiURL: apiURL,
		token:  token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendMessage отправляет сообщение в MAX
func (m *MAXClient) SendMessage(ctx context.Context, userID, text string) error {
	type sendMessageRequest struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}

	reqBody := sendMessageRequest{
		ChatID: userID,
		Text:   text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		m.apiURL+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", m.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("MAX API returned error: status %d", resp.StatusCode)
	}

	return nil
}

// ReceiveMessage получает сообщение (для Long Polling)
func (m *MAXClient) ReceiveMessage(ctx context.Context) (*interfaces.IncomingMessage, error) {
	req, err := http.NewRequestWithContext(ctx, "GET",
		m.apiURL+"/messages/updates", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", m.token)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to receive updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MAX API returned error: status %d", resp.StatusCode)
	}

	var incoming struct {
		UserID    string `json:"user_id"`
		MessageID string `json:"message_id"`
		Text      string `json:"text"`
		ChatID    string `json:"chat_id"`
		Timestamp int64  `json:"timestamp"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&incoming); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &interfaces.IncomingMessage{
		UserID:    incoming.UserID,
		MessageID: incoming.MessageID,
		Text:      incoming.Text,
		ChatID:    incoming.ChatID,
		Timestamp: incoming.Timestamp,
	}, nil
}

// SetWebhook устанавливает webhook для получения сообщений
func (m *MAXClient) SetWebhook(ctx context.Context, webhookURL string) error {
	type setWebhookRequest struct {
		URL string `json:"url"`
	}

	reqBody := setWebhookRequest{URL: webhookURL}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		m.apiURL+"/webhook", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", m.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("MAX API returned error: status %d", resp.StatusCode)
	}

	m.webhookURL = webhookURL
	return nil
}
