package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"socialmedia/chat/domain"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrDuplicateRequest = errors.New("duplicate follow request")
	ErrDuplicateFollow  = errors.New("duplicate follow relationship")
)

type Repository struct {
	db *sql.DB
}

func NewRepository(connString string) (*Repository, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully")

	if err := initDB(db); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	repo := &Repository{db: db}

	return repo, nil
}

func (r *Repository) CreateUser(ctx context.Context, id, username string) error {

	query := `
		INSERT INTO users_cache (id, username)
		VALUES ($1, $2)
		 ON CONFLICT (id) DO UPDATE
        SET username = EXCLUDED.username,
            updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, id, username)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *Repository) CreateConversation(ctx context.Context, isGroup bool, name string, userIDs []uuid.UUID) (*domain.Conversation, error) {
	query := `
		INSERT INTO conversations (is_group, name)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	var convo domain.Conversation
	convo.IsGroup = isGroup
	convo.Name = name

	err := r.db.QueryRowContext(ctx, query, isGroup, name).
		Scan(&convo.ID, &convo.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Katılımcıları ekle
	for _, uid := range userIDs {
		err := r.AddParticipant(ctx, convo.ID, uid)
		if err != nil {
			return nil, fmt.Errorf("failed to add participant: %w", err)
		}
	}

	return &convo, nil
}
func (r *Repository) AddParticipant(ctx context.Context, conversationID, userID uuid.UUID) error {
	query := `
		INSERT INTO conversation_participants (conversation_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to add participant: %w", err)
	}
	return nil
}
