package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"socialmedia/auth/domain"

	_ "github.com/lib/pq"
)

type PgRepository struct {
	db *sql.DB
}

func NewPgRepository(connString string) (*PgRepository, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	if err := InitDB(db); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &PgRepository{db: db}, nil
}
func InitDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		email VARCHAR(100) NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	log.Println("Users table created (if not exists)")
	return nil
}
func (r *PgRepository) SignUp(ctx context.Context, auth *domain.Auth) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, auth.Username, auth.Email, auth.Password)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}
