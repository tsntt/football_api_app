package model

import "context"

type Fan struct {
	ID               int    `json:"id" db:"id"`
	UserID           int    `json:"user_id" db:"user_id" validate:"required"`
	TeamID           int    `json:"team_id" db:"team_id" validate:"required"`
	NotificationType string `json:"notification_type" db:"notification_type"`
	Address          string `json:"address" db:"address"`
}

type IFanRepository interface {
	Create(ctx context.Context, fan *Fan) error
	GetAll(ctx context.Context) ([]Fan, error)
	GetByTeamID(ctx context.Context, teamID int) ([]Fan, error)
	GetByUserID(ctx context.Context, userID int) ([]Fan, error)
	DeleteByUserIDAndTeam(ctx context.Context, userID int, team string) error
}
