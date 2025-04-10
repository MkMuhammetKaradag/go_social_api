package postgres

import (
	"context"
	"fmt"
	"socialmedia/auth/domain"

	"golang.org/x/crypto/bcrypt"
)

func (r *Repository) hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func (r *Repository) SignUp(ctx context.Context, auth *domain.Auth) error {
	hashedPassword, err := r.hashPassword(auth.Password)
	if err != nil {
		return fmt.Errorf("hashing error: %w", err)
	}

	query := `INSERT INTO users (username, email, password, activation_code, activation_expiry) 
	          VALUES ($1, $2, $3, $4, $5)`

	_, err = r.db.ExecContext(ctx, query,
		auth.Username,
		auth.Email,
		hashedPassword,
		auth.ActivationCode,
		auth.ActivationExpiry,
	)

	if err != nil {
		return fmt.Errorf("insert error: %w", err)
	}

	return nil
}


