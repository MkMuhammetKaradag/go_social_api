package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"socialmedia/auth/domain"
	"time"
)

func (r *Repository) Activate(ctx context.Context, email, code string) (*domain.Auth, error) {
	const query = `
		SELECT id, username, email, activation_expiry
		FROM users 
		WHERE email = $1 AND activation_code = $2 AND is_active = false`

	var auth domain.Auth
	var expiry time.Time

	err := r.db.QueryRowContext(ctx, query, email, code).Scan(
		&auth.ID,
		&auth.Username,
		&auth.Email,
		&expiry,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("query error: %w", err)
	}

	if time.Now().After(expiry) {
		if _, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", auth.ID); err != nil {
			log.Printf("failed to delete expired user: %v", err)
		}
		return nil, ErrActivationExpired
	}

	const updateQuery = `
		UPDATE users 
		SET is_active = true, activation_code = NULL, activation_expiry = NULL
		WHERE id = $1
		RETURNING id, username, email`

	if err := r.db.QueryRowContext(ctx, updateQuery, auth.ID).Scan(
		&auth.ID,
		&auth.Username,
		&auth.Email,
	); err != nil {
		return nil, fmt.Errorf("update error: %w", err)
	}

	return &auth, nil
}
