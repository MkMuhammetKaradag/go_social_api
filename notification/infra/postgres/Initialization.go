package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	createUsersCacheTable = `
	CREATE TABLE IF NOT EXISTS users_cache (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(50) UNIQUE NOT NULL,
  avatar_url TEXT,
  is_private BOOLEAN DEFAULT false,
  updated_at TIMESTAMP DEFAULT NOW())`
)

func initDB(db *sql.DB) error {
	if _, err := db.Exec(createUsersCacheTable); err != nil {
		return fmt.Errorf("failed to create users_cache table: %w", err)
	}
	log.Println("Database tables initialized")
	return nil
}
