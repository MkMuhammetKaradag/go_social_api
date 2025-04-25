package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
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
func (r *Repository) UpdateUser(ctx context.Context, userID uuid.UUID, userName, avatarURL *string, isPrivate *bool) error {
	setClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	if userName != nil {
		setClauses = append(setClauses, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, *userName)
		argIndex++
	}
	if avatarURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("avatar_url = $%d", argIndex))
		args = append(args, *avatarURL)
		argIndex++
	}
	if isPrivate != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_private = $%d", argIndex))
		args = append(args, *isPrivate)
		argIndex++
	}

	if len(setClauses) == 0 {
		// Hiçbir alan güncellenmeyecekse işlemi boşuna yapma
		return nil
	}

	// updated_at her durumda güncellenir
	setClauses = append(setClauses, fmt.Sprintf("updated_at = NOW()"))

	query := fmt.Sprintf(`
		UPDATE users_cache
		SET %s
		WHERE id = $%d
	`, strings.Join(setClauses, ", "), argIndex)

	args = append(args, userID)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
