// handler/webhook_handler.go
package interfaces
type MessengerWebhookHandler interface {
    HandleIncomingMessage(message string, userID string) (string, error)
}