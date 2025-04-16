package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrDuplicateRequest = errors.New("duplicate follow request")
	ErrDuplicateFollow  = errors.New("duplicate follow relationship")
)

type Repository struct {
	db *sql.DB
}

func NewRepository(connString string) (*Repository, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully")

	if err := initDB(db); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	repo := &Repository{db: db}

	return repo, nil
}

func (r *Repository) CreateUser(ctx context.Context, id, username string) error {

	query := `
		INSERT INTO users_cache (id, username)
		VALUES ($1, $2)
		 ON CONFLICT (id) DO UPDATE
        SET username = EXCLUDED.username,
            updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, id, username)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

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
