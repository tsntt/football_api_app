package data

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/tsntt/footballapi/internal/model"
)

type FanRepository struct {
	db *sqlx.DB
}

func NewFanRepository(db *sqlx.DB) *FanRepository {
	return &FanRepository{db: db}
}

func (r *FanRepository) Create(ctx context.Context, fan *model.Fan) error {
	query := `INSERT INTO fans (user_id, team_id) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, fan.UserID, fan.TeamID).Scan(&fan.ID)
	if err != nil {
		return fmt.Errorf("failed to create fan: %w", err)
	}
	return nil
}

func (r *FanRepository) GetByTeamID(ctx context.Context, teamID int) ([]model.Fan, error) {
	fans := []model.Fan{}
	query := `SELECT id, user_id, team_id FROM fans WHERE team_id = $1`

	err := r.db.SelectContext(ctx, &fans, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fans by team: %w", err)
	}

	return fans, nil
}

func (r *FanRepository) GetByUserID(ctx context.Context, userID int) ([]model.Fan, error) {
	fans := []model.Fan{}
	query := `SELECT id, user_id, team_id FROM fans WHERE user_id = $1`

	err := r.db.SelectContext(ctx, &fans, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fans by user id: %w", err)
	}

	return fans, nil
}

func (r *FanRepository) DeleteByUserIDAndTeam(ctx context.Context, userID int, team string) error {
	query := `DELETE FROM fans WHERE user_id = $1 AND team = $2`

	result, err := r.db.ExecContext(ctx, query, userID, team)
	if err != nil {
		return fmt.Errorf("failed to delete fan subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}
