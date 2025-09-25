package model

import "context"

type Fan struct {
	ID     int    `json:"id" db:"id"`
	UserID int    `json:"user_id" db:"user_id" validate:"required"`
	Team   string `json:"team" db:"team" validate:"required"`
}

type IFanRepository interface {
	Create(ctx context.Context, fan *Fan) error
	GetByTeam(ctx context.Context, team string) ([]Fan, error)
	GetByUserID(ctx context.Context, userID int) ([]Fan, error)
}
