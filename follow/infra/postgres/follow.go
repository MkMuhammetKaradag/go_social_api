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
	checkQuery := `
        SELECT 1 FROM follow_requests 
        WHERE requester_id = $1 AND target_id = $2 AND status = 'pending'
    `

	var exists bool
	err := r.db.QueryRowContext(ctx, checkQuery, requesterID, targetID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check for existing request: %w", err)
	}

	// If a pending request exists, return error
	if exists {
		return ErrDuplicateRequest
	}

	// If no pending request exists, create a new one
	insertQuery := `
        INSERT INTO follow_requests (requester_id, target_id, status)
        VALUES ($1, $2, 'pending')
    `

	_, err = r.db.ExecContext(ctx, insertQuery, requesterID, targetID)
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

func (r *Repository) AcceptFollowRequest(ctx context.Context, requestID, currentUserID uuid.UUID) (uuid.UUID, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	// 1. follow_requests içinden requester ve target ID'yi çek
	var requesterID, targetID uuid.UUID
	selectQuery := `
		SELECT requester_id, target_id
		FROM follow_requests
		WHERE id = $1 AND status = 'pending'
	`
	err = tx.QueryRowContext(ctx, selectQuery, requestID).Scan(&requesterID, &targetID)
	if err == sql.ErrNoRows {
		return uuid.UUID{}, fmt.Errorf("follow request not found or already handled")
	}
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to fetch follow request: %w", err)
	}

	// 2. İsteğin gerçekten bu kullanıcıya gelip gelmediğini kontrol et
	if targetID != currentUserID {
		return uuid.UUID{}, fmt.Errorf("unauthorized: follow request not directed to this user")
	}

	// 3. follow_requests tablosunu güncelle (status = accepted)
	updateQuery := `
		UPDATE follow_requests 
		SET status = 'accepted'
		WHERE id = $1
	`
	_, err = tx.ExecContext(ctx, updateQuery, requestID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to update follow request: %w", err)
	}

	// 4. follows tablosuna takip ilişkisini yaz
	insertQuery := `
		INSERT INTO follows (follower_id, following_id)
		VALUES ($1, $2)
		ON CONFLICT (follower_id, following_id) DO NOTHING
	`
	_, err = tx.ExecContext(ctx, insertQuery, requesterID, targetID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert follow relation: %w", err)
	}

	return requesterID, nil
}

func (r *Repository) RejectFollowRequest(ctx context.Context, requestID, currentUserID uuid.UUID) (uuid.UUID, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	// 1. follow_requests içinden requester ve target ID'yi çek
	var requesterID, targetID uuid.UUID
	selectQuery := `
		SELECT requester_id, target_id
		FROM follow_requests
		WHERE id = $1 AND status = 'pending'
	`
	err = tx.QueryRowContext(ctx, selectQuery, requestID).Scan(&requesterID, &targetID)
	if err == sql.ErrNoRows {
		return uuid.UUID{}, fmt.Errorf("follow request not found or already handled")
	}
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to fetch follow request: %w", err)
	}

	// 2. İsteğin gerçekten bu kullanıcıya gelip gelmediğini kontrol et
	if targetID != currentUserID {
		return uuid.UUID{}, fmt.Errorf("unauthorized: follow request not directed to this user")
	}

	// 3. follow_requests tablosunu güncelle (status = accepted)
	updateQuery := `
		UPDATE follow_requests 
		SET status = 'rejectted'
		WHERE id = $1
	`
	_, err = tx.ExecContext(ctx, updateQuery, requestID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to update follow request: %w", err)
	}

	return requesterID, nil
}
