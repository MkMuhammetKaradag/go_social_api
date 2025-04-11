package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"socialmedia/user/domain"

	"github.com/google/uuid"
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

func (r *Repository) GetUserProfile(ctx context.Context, identifier string) (*domain.User, error) {
	var query string
	var row *sql.Row

	// Eğer UUID mi kontrol et
	if _, err := uuid.Parse(identifier); err == nil {
		// UUID ise id ile çek
		query = `SELECT id, username, email, bio, avatar_url, is_private, created_at, updated_at FROM users WHERE id = $1`
		row = r.db.QueryRowContext(ctx, query, identifier)
	} else {
		// Değilse username/email ile çek
		query = `SELECT id, username, email, bio, avatar_url, is_private, created_at, updated_at FROM users WHERE username = $1 OR email = $1`
		row = r.db.QueryRowContext(ctx, query, identifier)
	}
	// const query = `
	// 	SELECT
	// 		id,
	// 		username,
	// 		email,
	// 		bio,
	// 		avatar_url,
	// 		is_private,
	// 		created_at,
	// 	 	updated_at
	// 	FROM users
	// 	WHERE (username = $1 OR email = $1 OR id = $1)`

	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Bio,
		&user.AvatarURL,
		&user.IsPrivate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("query error: %w", err)
	}
	fmt.Println("repository içi", user)

	return &user, nil
}

func (r *Repository) CreateUser(ctx context.Context, id, username, email string) error {
	
	query := `
		INSERT INTO users (id, username, email)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, id, username, email)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
