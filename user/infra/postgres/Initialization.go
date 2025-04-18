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
	if _, err := db.Exec(createUsersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	if _, err := db.Exec(fallowTable); err != nil {
		return fmt.Errorf("failed to create fallow_cache table: %w", err)
	}

	if _, err := db.Exec(createBlockTable); err != nil {
		return fmt.Errorf("failed to create blocks_cache table: %w", err)
	}
	log.Println("Database tables initialized")
	return nil
}
