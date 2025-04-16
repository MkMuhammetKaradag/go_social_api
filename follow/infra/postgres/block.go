package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (r *Repository) BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Engelleme kaydı oluştur
	blockQuery := `INSERT INTO blocks (blocker_id, blocked_id) VALUES ($1, $2)`
	_, err = tx.ExecContext(ctx, blockQuery, blockerID, blockedID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("user already blocked")
		}
		return fmt.Errorf("failed to block user: %w", err)
	}

	// 2. Varsa follow ilişkisini sil
	deleteFollowQuery := `DELETE FROM follows WHERE 
        (follower_id = $1 AND following_id = $2) OR
        (follower_id = $2 AND following_id = $1)`
	_, err = tx.ExecContext(ctx, deleteFollowQuery, blockerID, blockedID)
	if err != nil {
		return fmt.Errorf("failed to remove follow relationship: %w", err)
	}

	// 3. Varsa follow isteklerini sil
	deleteRequestQuery := `DELETE FROM follow_requests WHERE 
        (requester_id = $1 AND target_id = $2) OR 
        (requester_id = $2 AND target_id = $1)`
	_, err = tx.ExecContext(ctx, deleteRequestQuery, blockerID, blockedID)
	if err != nil {
		return fmt.Errorf("failed to remove follow requests: %w", err)
	}

	// İşlemi tamamla
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Engeli kaldırma
func (r *Repository) UnblockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	query := `DELETE FROM blocks WHERE blocker_id = $1 AND blocked_id = $2`
	result, err := r.db.ExecContext(ctx, query, blockerID, blockedID)
	if err != nil {
		return fmt.Errorf("failed to unblock user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("block relationship not found")
	}

	return nil
}

// Kullanıcının engellenip engellenmediğini kontrol etme
func (r *Repository) IsBlocked(ctx context.Context, userID, targetID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM blocks WHERE blocker_id = $1 AND blocked_id = $2)`

	var isBlocked bool
	err := r.db.QueryRowContext(ctx, query, userID, targetID).Scan(&isBlocked)
	if err != nil {
		return false, fmt.Errorf("failed to check block status: %w", err)
	}

	return isBlocked, nil
}

// Herhangi birinin engellemiş olmasını kontrol etme (çift yönlü)
func (r *Repository) HasBlockRelationship(ctx context.Context, userID1, userID2 uuid.UUID) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM blocks 
            WHERE (blocker_id = $1 AND blocked_id = $2) OR 
                  (blocker_id = $2 AND blocked_id = $1)
        )
    `

	var hasBlock bool
	err := r.db.QueryRowContext(ctx, query, userID1, userID2).Scan(&hasBlock)
	if err != nil {
		return false, fmt.Errorf("failed to check block relationship: %w", err)
	}

	return hasBlock, nil
}
