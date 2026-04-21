package messenger

type Client interface {
	SendMessage(userID string, text string) error
	SetWebhook(url string) error
}

type maxClient struct {
	apiURL string
	token  string
}
