package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	createUsersTable = `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		username VARCHAR(50) NOT NULL UNIQUE,
		email VARCHAR(100) NOT NULL UNIQUE,
		password TEXT NOT NULL,
		is_active BOOLEAN DEFAULT false,
		activation_code VARCHAR(4),
		activation_expiry TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)` 
	// id SERIAL PRIMARY KEY,

	createForgotPasswordsTable = `
	CREATE TABLE IF NOT EXISTS forgot_passwords (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		token TEXT NOT NULL,
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)` // id SERIAL PRIMARY KEY,
)

func initDB(db *sql.DB) error {
	if _, err := db.Exec(createUsersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	if _, err := db.Exec(createForgotPasswordsTable); err != nil {
		return fmt.Errorf("failed to create forgot_passwords table: %w", err)
	}

	log.Println("Database tables initialized")
	return nil
}
