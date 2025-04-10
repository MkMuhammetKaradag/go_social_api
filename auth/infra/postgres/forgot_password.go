package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"socialmedia/auth/domain"
)

func (r *Repository) RequestForgotPassword(ctx context.Context, fp *domain.ForgotPassword) (string, error) {
	const userQuery = `SELECT id, username FROM users WHERE email = $1`

	var userID int
	var username string

	if err := r.db.QueryRowContext(ctx, userQuery, fp.Email).Scan(&userID, &username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("query error: %w", err)
	}

	const insertQuery = `INSERT INTO forgot_passwords (user_id, token, expires_at) VALUES ($1, $2, $3)`
	if _, err := r.db.ExecContext(ctx, insertQuery, userID, fp.Token, fp.ExpiresAt); err != nil {
		return "", fmt.Errorf("insert error: %w", err)
	}

	return username, nil
}

