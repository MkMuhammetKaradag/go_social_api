package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	ErrUserNotFound = errors.New("user not found")
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
