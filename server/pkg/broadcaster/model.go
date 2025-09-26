package broadcaster

import (
	"context"
)

type NotificationTarget struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	Address  string         `json:"address"`
	Metadata map[string]int `json:"metadata"`
}

type IBroadcastService interface {
	SendNotification(ctx context.Context, message string, targets []NotificationTarget) error
	SendNotificationWithID(ctx context.Context, notificationID, message string, targets []NotificationTarget) error
}
