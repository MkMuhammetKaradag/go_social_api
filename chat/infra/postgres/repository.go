package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"socialmedia/chat/domain"
	"strings"

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

func (r *Repository) UpdateUser(ctx context.Context, userID uuid.UUID, userName, avatarURL *string, isPrivate *bool) error {
	setClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	if userName != nil {
		setClauses = append(setClauses, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, *userName)
		argIndex++
	}
	if avatarURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("avatar_url = $%d", argIndex))
		args = append(args, *avatarURL)
		argIndex++
	}
	if isPrivate != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_private = $%d", argIndex))
		args = append(args, *isPrivate)
		argIndex++
	}

	if len(setClauses) == 0 {
		// Hiçbir alan güncellenmeyecekse işlemi boşuna yapma
		return nil
	}

	// updated_at her durumda güncellenir
	setClauses = append(setClauses, fmt.Sprintf("updated_at = NOW()"))

	query := fmt.Sprintf(`
		UPDATE users_cache
		SET %s
		WHERE id = $%d
	`, strings.Join(setClauses, ", "), argIndex)

	args = append(args, userID)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *Repository) CreateConversation(ctx context.Context, currrentUserID uuid.UUID, isGroup bool, name string, userIDs []uuid.UUID) (*domain.Conversation, *[]domain.BlockedParticipant, error) {
	// Bloklanmış kullanıcıları tutacak slice
	var blockedParticipants []domain.BlockedParticipant
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
		return nil, nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Tüm kullanıcılar arasındaki blok ilişkilerini kontrol et
	for i, uid1 := range userIDs {
		// if uid1 == currrentUserID {
		// 	continue
		// }
		for j, uid2 := range userIDs {
			if i != j {
				blocked, err := r.IsBlocked(ctx, uid1, uid2)
				if err != nil {
					return nil, nil, err
				}
				if blocked {
					blockedParticipants = append(blockedParticipants, domain.BlockedParticipant{
						BlockerID: uid1,
						BlockedID: uid2,
					})
				}
			}
		}
	}

	// Katılımcıları ekle
	for _, uid := range userIDs {
		err := r.AddParticipant(ctx, convo.ID, uid, currrentUserID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to add participant: %w", err)
		}
	}

	return &convo, &blockedParticipants, nil
}

func (r *Repository) GetBlockedParticipantsForUser(ctx context.Context, conversationID, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `
        SELECT b.blocked_id
        FROM blocks_cache b
        INNER JOIN conversation_participants cp1 ON b.blocker_id = $1
        INNER JOIN conversation_participants cp2 ON b.blocked_id = cp2.user_id
        WHERE cp1.conversation_id = $2 AND cp2.conversation_id = $2
    `

	var blockedUsers []uuid.UUID
	rows, err := r.db.QueryContext(ctx, query, userID, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var blockedID uuid.UUID
		if err := rows.Scan(&blockedID); err != nil {
			return nil, err
		}
		blockedUsers = append(blockedUsers, blockedID)
	}

	return blockedUsers, nil
}
func (r *Repository) AddParticipant(ctx context.Context, conversationID, userID uuid.UUID, currentUserID uuid.UUID) error {
	isPrivate, err := r.IsUserPrivate(ctx, userID)
	if err != nil {
		return err
	}
	// blocked, err := r.HasBlockRelationship(ctx, currentUserID, userID)
	// if err != nil {
	// 	return err
	// }
	// if blocked {
	// 	return nil
	// }

	if isPrivate {
		areFriends, err := r.AreUsersFriends(ctx, currentUserID, userID)
		if err != nil {
			return err
		}
		if !areFriends {
			return nil //fmt.Errorf("user %s is private and not friends with %s", userID, currentUserID)
		}
	}

	query := `
		INSERT INTO conversation_participants (conversation_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err = r.db.ExecContext(ctx, query, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to add participant: %w", err)
	}
	return nil
}

func (r *Repository) CreateMessage(ctx context.Context, conversationID, senderID uuid.UUID, content string, attachmentURLs []string, attachmentTypes []string) (*domain.Message, error) {
	// First, check if the sender is a participant in the conversation
	isParticipant, err := r.IsParticipant(ctx, conversationID, senderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check participant status: %w", err)
	}

	if !isParticipant {
		return nil, fmt.Errorf("user is not a participant in this conversation")
	}

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert the message
	query := `
        INSERT INTO messages (conversation_id, sender_id, content)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, is_edited, deleted_at
    `

	var msg domain.Message
	msg.ConversationID = conversationID
	msg.UserID = senderID
	msg.Content = content

	err = tx.QueryRowContext(ctx, query, conversationID, senderID, content).
		Scan(&msg.ID, &msg.CreatedAt, &msg.IsEdited, &msg.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Add attachments if provided
	if len(attachmentURLs) > 0 {
		for i, fileURL := range attachmentURLs {
			var fileType string
			if i < len(attachmentTypes) {
				fileType = attachmentTypes[i]
			}

			attachQuery := `
                INSERT INTO attachments (message_id, file_url, file_type)
                VALUES ($1, $2, $3)
                RETURNING id
            `

			var attachmentID uuid.UUID
			err = tx.QueryRowContext(ctx, attachQuery, msg.ID, fileURL, fileType).
				Scan(&attachmentID)
			if err != nil {
				return nil, fmt.Errorf("failed to add attachment: %w", err)
			}

			attachment := domain.Attachment{
				ID:        attachmentID,
				MessageID: msg.ID,
				FileURL:   fileURL,
				FileType:  fileType,
			}

			msg.Attachments = append(msg.Attachments, attachment)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &msg, nil
}

// Helper function to check if a user is a participant in a conversation
func (r *Repository) IsParticipant(ctx context.Context, conversationID, userID uuid.UUID) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1 FROM conversation_participants 
            WHERE conversation_id = $1 AND user_id = $2
        )
    `

	var isParticipant bool
	err := r.db.QueryRowContext(ctx, query, conversationID, userID).Scan(&isParticipant)
	if err != nil {
		return false, fmt.Errorf("failed to check participant status: %w", err)
	}

	return isParticipant, nil
}

func (r *Repository) IsUserPrivate(ctx context.Context, userID uuid.UUID) (bool, error) {
	query := `SELECT is_private FROM users_cache WHERE id = $1`
	var isPrivate bool
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&isPrivate)
	if err != nil {
		return false, fmt.Errorf("failed to check privacy status: %w", err)
	}
	return isPrivate, nil
}

func (r *Repository) AreUsersFriends(ctx context.Context, userA, userB uuid.UUID) (bool, error) {
	query := `
		SELECT COUNT(*) FROM follows_cache
		WHERE follower_id = $1 AND following_id = $2 AND status = 'following'
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, userA, userB).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check friendship: %w", err)
	}
	return count > 0, nil
}
func (r *Repository) GetParticipants(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error) {
	// Örnek SQL sorgusu
	query := `SELECT user_id FROM conversation_participants WHERE conversation_id = $1`

	rows, err := r.db.QueryContext(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		participants = append(participants, userID)
	}

	return participants, nil
}
