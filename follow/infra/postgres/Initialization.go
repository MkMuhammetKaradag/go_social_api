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

	createFallowTable = `
	CREATE TABLE IF NOT EXISTS follows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    follower_id UUID NOT NULL,
    following_id UUID NOT NULL,
    followed_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (follower_id, following_id),
    CHECK (follower_id != following_id)
)`

	reateFallowRequestTable = `
	CREATE TABLE IF NOT EXISTS follow_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requester_id UUID NOT NULL,
    target_id UUID NOT NULL,
    requested_at TIMESTAMP DEFAULT NOW(),
    status TEXT DEFAULT 'pending',  -- 'pending', 'accepted', 'rejected' 
	
    CHECK (requester_id != target_id)
)`
	// UNIQUE (requester_id, target_id),
	createBlockTable = `
	CREATE TABLE IF NOT EXISTS blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blocker_id UUID NOT NULL,
    blocked_id UUID NOT NULL,
    blocked_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (blocker_id, blocked_id),
    CHECK (blocker_id != blocked_id)
)`
)

func initDB(db *sql.DB) error {
	if _, err := db.Exec(createUsersCacheTable); err != nil {
		return fmt.Errorf("failed to create users_cache table: %w", err)
	}

	if _, err := db.Exec(createFallowTable); err != nil {
		return fmt.Errorf("failed to create fallow table: %w", err)
	}
	if _, err := db.Exec(reateFallowRequestTable); err != nil {
		return fmt.Errorf("failed to create fallow_requests table: %w", err)
	}

	if _, err := db.Exec(createBlockTable); err != nil {
		return fmt.Errorf("failed to create blocks table: %w", err)
	}
	log.Println("Database tables initialized")
	return nil
}
