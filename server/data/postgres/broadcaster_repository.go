package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/tsntt/footballapi/internal/model"
)

type BroadcastRepository struct {
	db *sqlx.DB
}

func NewBroadcastRepository(db *sqlx.DB) *BroadcastRepository {
	return &BroadcastRepository{db: db}
}

func (r *BroadcastRepository) Create(ctx context.Context, broadcast *model.BroadcastMessage) error {
	query := `
		INSERT INTO broadcasted_messages (match_id, message_content_hash, status) 
		VALUES ($1, $2, $3) 
		RETURNING id, sent_at`

	err := r.db.QueryRowContext(ctx, query, broadcast.MatchID, broadcast.MessageContentHash, broadcast.Status).
		Scan(&broadcast.ID, &broadcast.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create broadcast message: %w", err)
	}

	return nil
}

func (r *BroadcastRepository) GetByMatchID(ctx context.Context, matchID int) (*model.BroadcastMessage, error) {
	broadcast := &model.BroadcastMessage{}
	query := `SELECT id, match_id, message_content_hash, created_at, status FROM broadcasted_messages WHERE match_id = $1`

	err := r.db.GetContext(ctx, broadcast, query, matchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("broadcast message not found")
		}
		return nil, fmt.Errorf("failed to get broadcast message: %w", err)
	}

	return broadcast, nil
}

func (r *BroadcastRepository) Update(ctx context.Context, broadcast *model.BroadcastMessage) error {
	query := `UPDATE broadcasted_messages SET status = $2 WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, broadcast.Status, broadcast.ID)
	if err != nil {
		return fmt.Errorf("failed to update broadcast message: %w", err)
	}

	return nil
}
