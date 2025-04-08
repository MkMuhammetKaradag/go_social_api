package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"socialmedia/auth/domain"
	"time"

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
	queryUsers := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		email VARCHAR(100) NOT NULL UNIQUE,
		password TEXT NOT NULL,
		is_active BOOLEAN DEFAULT false,
		activation_code VARCHAR(4),
    	activation_expiry TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);` // created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

	queryForgotPasswords := `
	CREATE TABLE IF NOT EXISTS forgotPasswords (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id) ON DELETE CASCADE,
		token TEXT NOT NULL,
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);`

	_, err := db.Exec(queryUsers)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	_, err = db.Exec(queryForgotPasswords)
	if err != nil {
		return fmt.Errorf("failed to create resetPasswords table: %w", err)
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
	query := `INSERT INTO users (username, email, password,activation_code,activation_expiry) VALUES ($1, $2, $3,$4,$5)`

	_, err = r.db.ExecContext(ctx, query, auth.Username, auth.Email, hashedPassword, auth.ActivationCode, auth.ActivationExpiry)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (r *PgRepository) RequestForgotPassword(ctx context.Context, forgotPassword *domain.ForgotPassword) (*string, error) {
	// Kullanıcıyı e-posta adresine göre bulalım
	var userID int
	var username string
	queryUser := `SELECT id,username  FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, queryUser, forgotPassword.Email).Scan(&userID, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			// Eğer kullanıcı bulunamazsa, hata vermek yerine sadece işlem yapma
			return nil, fmt.Errorf("user not found with the provided email")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	queryResetPassword := `
		INSERT INTO forgotPasswords (user_id, token, expires_at) 
		VALUES ($1, $2, $3)`
	_, err = r.db.ExecContext(ctx, queryResetPassword, userID, forgotPassword.Token, forgotPassword.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert reset password record: %w", err)
	}

	return &username, nil
}

func (r *PgRepository) Activate(ctx context.Context, userEmail string, activationCode string) (*domain.Auth, error) {
	// 1. Kullanıcıyı email ve activation code ile bul
	query := `
SELECT id, username, email, activation_code, activation_expiry
FROM users 
WHERE email = $1 
AND activation_code = $2
AND is_active = false` // Sadece aktif olmayan kullanıcılar için

	var auth domain.Auth
	var expiryTime time.Time
	err := r.db.QueryRowContext(ctx, query, userEmail, activationCode).Scan(
		&auth.ID,
		&auth.Username,
		&auth.Email,
		&auth.ActivationCode,
		&expiryTime,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid activation code or email")
		}
		return nil, fmt.Errorf("failed to activate user: %w", err)
	}

	// 2. Aktivasyon süresini kontrol et
	if time.Now().After(expiryTime) {
		// Süre dolmuşsa kullanıcıyı sil (opsiyonel)
		_, _ = r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", auth.ID)
		return nil, fmt.Errorf("activation code expired")
	}

	// 3. Kullanıcıyı aktif et
	updateQuery := `
UPDATE users 
SET is_active = true, 
	activation_code = NULL,
	activation_expiry = NULL
WHERE id = $1
RETURNING id, username, email`

	// Aktif edilmiş halini döndür (password/code olmadan)
	err = r.db.QueryRowContext(ctx, updateQuery, auth.ID).Scan(
		&auth.ID,
		&auth.Username,
		&auth.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user status: %w", err)
	}

	// 4. Hassas alanları temizle
	auth.Password = ""
	auth.ActivationCode = ""

	return &auth, nil
}
func (r *PgRepository) SignIn(ctx context.Context, identifier, password string) (*domain.Auth, error) {
	query := `
		SELECT id, username, email, password 
        FROM users 
        WHERE (username = $1 OR email = $1)
		AND is_active = true`

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

func (r *PgRepository) StartCleanupJob(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			// 5 dakikadan eski ve aktif edilmemiş kullanıcıları sil
			query := `
                DELETE FROM users 
                WHERE is_active = false 
                AND activation_expiry < NOW()`

			_, err := r.db.Exec(query)
			if err != nil {
				log.Printf("Cleanup job error: %v", err)
				continue
			}
			log.Printf("Cleanup job: Inactive users deleted")
		}
	}()
}
