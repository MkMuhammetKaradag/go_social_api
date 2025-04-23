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

	fallowTable = `
  CREATE TABLE IF NOT EXISTS follows_cache  (
follower_id UUID NOT NULL,
  following_id UUID NOT NULL,
  status TEXT NOT NULL, -- 'following', 'pending', 'none'
  PRIMARY KEY (follower_id, following_id),
  CHECK (follower_id != following_id)
)`

	createBlockTable = `
  CREATE TABLE IF NOT EXISTS blocks_cache (
  blocker_id UUID NOT NULL,
  blocked_id UUID NOT NULL,
  PRIMARY KEY (blocker_id, blocked_id),
  CHECK (blocker_id != blocked_id)
)`
)

func initDB(db *sql.DB) error {
	if _, err := db.Exec(createUsersCacheTable); err != nil {
		return fmt.Errorf("failed to create users_cache table: %w", err)
	}

	if _, err := db.Exec(fallowTable); err != nil {
		return fmt.Errorf("failed to create fallow table: %w", err)
	}
	if _, err := db.Exec(createBlockTable); err != nil {
		return fmt.Errorf("failed to create fallow_requests table: %w", err)
	}

	log.Println("Database tables initialized")
	return nil
}
