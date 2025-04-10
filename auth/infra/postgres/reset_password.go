package postgres

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (r *Repository) ResetPassword(ctx context.Context, token, password string) (*int, error) {

	userID, expiryTime, err := r.getUserIDFromToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("password reset failed: %w", err)
	}

	if time.Now().After(expiryTime) {
		if err := r.deleteExpiredToken(ctx, token); err != nil {
			return nil, fmt.Errorf("failed to delete expired token: %w", err)
		}
		return nil, fmt.Errorf("token has expired")
	}

	hashedPassword, err := r.hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	if err := r.updateUserPassword(ctx, userID, hashedPassword); err != nil {
		return nil, fmt.Errorf("failed to update password: %w", err)
	}

	if err := r.deleteToken(ctx, token); err != nil {

		log.Printf("warning: failed to delete token after password reset: %v", err)
	}

	return &userID, nil
}
