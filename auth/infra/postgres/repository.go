package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"socialmedia/auth/domain"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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
func (r *PgRepository) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (r *PgRepository) SignUp(ctx context.Context, auth *domain.Auth) error {
	hashedPassword, err := r.HashPassword(auth.Password) // Şifreyi hash'le
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`

	_, err = r.db.ExecContext(ctx, query, auth.Username, auth.Email, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (r *PgRepository) SignIn(ctx context.Context, identifier, password string) (*domain.Auth, error) {
	query := `SELECT id, username, email, password FROM users WHERE username = $1 OR email = $1`

	row := r.db.QueryRowContext(ctx, query, identifier)

	var auth domain.Auth
	var hashedPassword string
	err := row.Scan(&auth.ID, &auth.Username, &auth.Email, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &auth, nil
}
