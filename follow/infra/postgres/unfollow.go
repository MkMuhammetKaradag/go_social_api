package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)


func (r *Repository) DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `
		DELETE FROM follows
		WHERE follower_id = $1 AND following_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to delete follow relationship: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows 
	}

	return nil
}

func (r *Repository) DeleteFollowRequest(ctx context.Context, requesterID, targetID uuid.UUID) error {
	query := `
		DELETE FROM follow_requests
		WHERE requester_id = $1 AND target_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, requesterID, targetID)
	if err != nil {
		return fmt.Errorf("failed to delete follow request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows 
	}

	return nil
}
