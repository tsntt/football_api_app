package model

import (
	"context"
	"time"
)

type BroadcastMessage struct {
	ID                 int       `json:"id" db:"id"`
	MatchID            int       `json:"match_id" db:"match_id"`
	MessageContentHash string    `json:"message_content_hash" db:"message_content_hash"`
	Status             string    `json:"status" db:"status"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

type IBroadcastRepository interface {
	Create(ctx context.Context, broadcast *BroadcastMessage) error
	GetByMatchID(ctx context.Context, matchID int) (*BroadcastMessage, error)
	Update(ctx context.Context, broadcast *BroadcastMessage) error
}
