package postgres

import (
	"context"

	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) CreateFollow(ctx context.Context, followerID, followingID uuid.UUID, status string) error {
	query := `
    INSERT INTO follows_cache (follower_id, following_id, status)
    VALUES ($1, $2, $3)
    ON CONFLICT (follower_id, following_id) 
    DO UPDATE SET status = $3
`
	_, err := r.db.ExecContext(ctx, query, followerID, followingID, status)
	if err != nil {
		return fmt.Errorf("failed to create or update follow relationship: %w", err)
	}

	return nil
}

func (r *Repository) DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `
		DELETE FROM follows_cache
		WHERE follower_id = $1 AND following_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to delete follow relationship: %w", err)
	}



	return nil
}
