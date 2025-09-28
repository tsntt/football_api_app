package broadcast

import "context"

type NotificationType string

const (
	Email     NotificationType = "email"
	SMS       NotificationType = "sms"
	Push      NotificationType = "push"
	WebSocket NotificationType = "websocket"
	Webhook   NotificationType = "webhook"
)

type Channel struct {
	ID            int            `json:"id" db:"id"`
	Name          string         `json:"name" db:"name"`
	Subscriptions []Subscription `json:"subscriptions" db:"subscriptions"`
}

type Subscription struct {
	UserID           int              `json:"user_id" db:"user_id"`
	ChannelID        int              `json:"channel_id" db:"channel_id"`
	NotificationType NotificationType `json:"notification_type" db:"notification_type"`
	Address          string           `json:"address" db:"address"`
}

type Message struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type BroadcastStatus struct {
	ChannelID    int      `json:"channel_id"`
	TotalToSend  int      `json:"total_sent"`
	SentCount    int      `json:"sent_count"`
	FailedCount  int      `json:"failed_count"`
	IsCompleted  bool     `json:"is_completed"`
	ErrorDetails []string `json:"error_details"`
}

type IBroadcaster interface {
	Send(ctx context.Context, subscription Subscription, message Message) error
}
