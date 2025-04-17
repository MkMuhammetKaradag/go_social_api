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

func (r *Repository) DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `
		DELETE FROM follows_cache
		WHERE follower_id = $1 AND following_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to delete follow relationship: %w", err)
	}

	// rowsAffected, err := _.RowsAffected()
	// if err != nil {
	// 	return fmt.Errorf("failed to get affected rows: %w", err)
	// }

	// if rowsAffected == 0 {
	// 	return sql.ErrNoRows
	// }

	return nil
}
