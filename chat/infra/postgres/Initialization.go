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

	conversationTable = `
  CREATE TABLE IF NOT EXISTS conversations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  is_group BOOLEAN DEFAULT false,
  name TEXT,
  created_at TIMESTAMP DEFAULT NOW()
)`

	conversationParticipantTable = `
CREATE TABLE IF NOT EXISTS conversation_participants (
  conversation_id UUID NOT NULL,
  user_id UUID NOT NULL,
  joined_at TIMESTAMP DEFAULT NOW(),
  is_admin BOOLEAN DEFAULT false,
  PRIMARY KEY (conversation_id, user_id),
  FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
)`

	messagesTable = `
CREATE TABLE IF NOT EXISTS messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  conversation_id UUID NOT NULL,
  sender_id UUID NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  is_edited BOOLEAN DEFAULT false,
  deleted_at TIMESTAMP,
  FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
)`

	attachmentTable = `
CREATE TABLE IF NOT EXISTS attachments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  message_id UUID NOT NULL,
  file_url TEXT NOT NULL,
  file_type TEXT,
  FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE
)`
	messagesReadTable = `
  CREATE TABLE IF NOT EXISTS message_reads (
  message_id UUID NOT NULL,
  user_id UUID NOT NULL,
  read_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (message_id, user_id),
  FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE
)
`
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
	if _, err := db.Exec(conversationTable); err != nil {
		return fmt.Errorf("failed to create fallow_requests table: %w", err)
	}
	if _, err := db.Exec(conversationParticipantTable); err != nil {
		return fmt.Errorf("failed to create fallow_requests table: %w", err)
	}
	if _, err := db.Exec(messagesTable); err != nil {
		return fmt.Errorf("failed to create fallow_requests table: %w", err)
	}

	if _, err := db.Exec(attachmentTable); err != nil {
		return fmt.Errorf("failed to create fallow_requests table: %w", err)
	}
	if _, err := db.Exec(messagesReadTable); err != nil {
		return fmt.Errorf("failed to create message_reads table: %w", err)
	}
	log.Println("Database tables initialized")
	return nil
}
