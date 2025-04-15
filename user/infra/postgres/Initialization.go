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
  username VARCHAR(50) UNIQUE NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  bio TEXT,
  avatar_url TEXT,
  banner_url TEXT,
   location VARCHAR(100),
   website VARCHAR(200),
  is_private BOOLEAN DEFAULT false,
   is_verified BOOLEAN DEFAULT false,
   followers_count INTEGER DEFAULT 0,
  following_count INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`
)

func initDB(db *sql.DB) error {
	if _, err := db.Exec(createUsersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	log.Println("Database tables initialized")
	return nil
}
