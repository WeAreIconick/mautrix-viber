package viber

// Incoming webhook payload minimal structures

// Event represents a Viber webhook event type.
type Event string

const (
	// EventMessage is a message event.
	EventMessage Event = "message"
	// EventSubscribed is a subscription event.
	EventSubscribed = "subscribed"
	// EventUnsubscribed is an unsubscription event.
	EventUnsubscribed = "unsubscribed"
	// EventConversation is a conversation start event.
	EventConversation = "conversation_started"
)

// Sender represents a message sender.
type Sender struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar,omitempty"`
	Language string `json:"language,omitempty"`
	Country  string `json:"country,omitempty"`
}

// Message represents a Viber message.
type Message struct {
	Type      string `json:"type"`
	Text      string `json:"text,omitempty"`
	Media     string `json:"media,omitempty"`
	FileName  string `json:"file_name,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	// Group or chat identifiers (may vary by Viber API version)
	ChatID string `json:"chat_id,omitempty"`
}

// WebhookRequest represents an incoming Viber webhook request.
type WebhookRequest struct {
	Event   Event   `json:"event"`
	Sender  Sender  `json:"sender"`
	Message Message `json:"message"`
	// Token/hostname fields commonly present in Viber webhooks
	MessageToken int64  `json:"message_token,omitempty"`
	ChatHostname string `json:"chat_hostname,omitempty"`
}

// WebhookResponse represents a Viber webhook set response.
type WebhookResponse struct {
	Status        int    `json:"status"`
	StatusMessage string `json:"status_message"`
}
