package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func (r *Repository) startCleanupJob(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		if err := r.cleanupExpiredActivations(); err != nil {
			log.Printf("cleanup error: %v", err)
		}
	}
}

func (r *Repository) cleanupExpiredActivations() error {
	const query = `DELETE FROM users WHERE is_active = false AND activation_expiry < NOW()`
	_, err := r.db.Exec(query)
	return err
}

func (r *Repository) getUserIDFromToken(ctx context.Context, token string) (int, time.Time, error) {
	var userID int
	var expiryTime time.Time

	query := `SELECT user_id, expires_at FROM forgot_passwords WHERE token = $1`
	err := r.db.QueryRowContext(ctx, query, token).Scan(&userID, &expiryTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, time.Time{}, fmt.Errorf("token not found")
		}
		return 0, time.Time{}, fmt.Errorf("database error: %w", err)
	}

	return userID, expiryTime, nil
}

func (r *Repository) deleteExpiredToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM forgot_passwords WHERE token = $1", token)
	if err != nil {
		return fmt.Errorf("failed to delete expired token: %w", err)
	}
	return nil
}

func (r *Repository) updateUserPassword(ctx context.Context, userID int, hashedPassword string) error {
	updateQuery := `UPDATE users SET password = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, updateQuery, hashedPassword, userID)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *Repository) deleteToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM forgot_passwords WHERE token = $1", token)
	return err
}

