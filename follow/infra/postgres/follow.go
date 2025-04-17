package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (r *Repository) IsPrivate(ctx context.Context, userID uuid.UUID) (bool, error) {
	query := `SELECT is_private FROM users_cache WHERE id = $1`

	var isPrivate bool
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&isPrivate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("user not found: %w", err)
		}
		return false, fmt.Errorf("failed to query user privacy status: %w", err)
	}

	return isPrivate, nil
}
func (r *Repository) CreateFollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `
		INSERT INTO follows (follower_id, following_id)
		VALUES ($1, $2)
	`

	_, err := r.db.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		// Check for unique constraint violation (already following)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrDuplicateFollow
		}
		return fmt.Errorf("failed to create follow relationship: %w", err)
	}

	return nil
}

func (r *Repository) CreateFollowRequest(ctx context.Context, requesterID, targetID uuid.UUID) error {
	query := `
		INSERT INTO follow_requests (requester_id, target_id, status)
		VALUES ($1, $2, 'pending')
	`

	_, err := r.db.ExecContext(ctx, query, requesterID, targetID)
	if err != nil {
		// Check for unique constraint violation (request already exists)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrDuplicateRequest
		}
		return fmt.Errorf("failed to create follow request: %w", err)
	}

	return nil
}
func (r *Repository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM follows
			WHERE follower_id = $1 AND following_id = $2
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, followerID, followingID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check follow relationship: %w", err)
	}

	return exists, nil
}
