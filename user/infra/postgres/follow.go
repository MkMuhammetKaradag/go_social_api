package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (r *Repository) CreateFollow(ctx context.Context, followerID, followingID uuid.UUID, status string) error {
	query := `
		INSERT INTO follows_cache (follower_id, following_id,status)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(ctx, query, followerID, followingID, status)
	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("failed to duplicate ")
		}
		return fmt.Errorf("failed to create follow relationship: %w", err)
	}

	return nil
}
