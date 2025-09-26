package model

import (
	"context"
	"time"
)

type BroadcastMessage struct {
	ID      int       `json:"id" db:"id"`
	MatchID int       `json:"match_id" db:"match_id"`
	Message string    `json:"message" db:"message"`
	SentAt  time.Time `json:"sent_at" db:"sent_at"`
	Status  string    `json:"status" db:"status"`
}

type IBroadcastRepository interface {
	Create(ctx context.Context, broadcast *BroadcastMessage) error
	GetByMatchID(ctx context.Context, matchID int) (*BroadcastMessage, error)
	Update(ctx context.Context, broadcast *BroadcastMessage) error
}
