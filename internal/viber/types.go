package viber

// Incoming webhook payload minimal structures

type Event string

const (
	EventMessage   Event = "message"
	EventSubscribed      = "subscribed"
	EventUnsubscribed    = "unsubscribed"
	EventConversation    = "conversation_started"
)

type Sender struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar,omitempty"`
	Language string `json:"language,omitempty"`
	Country  string `json:"country,omitempty"`
}

type Message struct {
	Type    string `json:"type"`
	Text    string `json:"text,omitempty"`
	Media   string `json:"media,omitempty"`
    FileName string `json:"file_name,omitempty"`
    Thumbnail string `json:"thumbnail,omitempty"`
    // Group or chat identifiers (may vary by Viber API version)
    ChatID   string `json:"chat_id,omitempty"`
}

type WebhookRequest struct {
	Event   Event   `json:"event"`
	Sender  Sender  `json:"sender"`
	Message Message `json:"message"`
    // Token/hostname fields commonly present in Viber webhooks
    MessageToken int64  `json:"message_token,omitempty"`
    ChatHostname string `json:"chat_hostname,omitempty"`
}

// Webhook set response

type WebhookResponse struct {
	Status        int    `json:"status"`
	StatusMessage string `json:"status_message"`
}
