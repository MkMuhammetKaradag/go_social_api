package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"socialmedia/auth/domain"

	"golang.org/x/crypto/bcrypt"
)

func (r *Repository) SignIn(ctx context.Context, identifier, password string) (*domain.Auth, error) {
	const query = `
		SELECT id, username, email, password 
		FROM users 
		WHERE (username = $1 OR email = $1) AND is_active = true`

	var auth domain.Auth
	var hashedPassword string

	err := r.db.QueryRowContext(ctx, query, identifier).Scan(
		&auth.ID,
		&auth.Username,
		&auth.Email,
		&hashedPassword,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("query error: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &auth, nil
}

