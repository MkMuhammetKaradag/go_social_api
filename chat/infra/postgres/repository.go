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

func (r *Repository) CreateConversation(ctx context.Context, currentUserID uuid.UUID, isGroup bool, name string, userIDs []uuid.UUID) (*domain.Conversation, *[]domain.BlockedParticipant, error) {
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
	err = r.addInitialAdmin(ctx, convo.ID, currentUserID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to add initial admin: %w", err)
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
		if uid == currentUserID {
			continue
		}
		err := r.AddParticipant(ctx, convo.ID, uid, currentUserID)
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
func (r *Repository) addInitialAdmin(ctx context.Context, conversationID, userID uuid.UUID) error {
	query := `
		INSERT INTO conversation_participants (conversation_id, user_id, is_admin)
		VALUES ($1, $2, true)
	`
	_, err := r.db.ExecContext(ctx, query, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to add initial admin: %w", err)
	}
	return nil
}
func (r *Repository) AddParticipant(ctx context.Context, conversationID, userID, addedByUserID uuid.UUID) error {

	isAdderAdmin, err := r.IsUserAdmin(ctx, conversationID, addedByUserID)
	if err != nil {
		return fmt.Errorf("failed to check admin rights: %w", err)
	}
	if !isAdderAdmin {
		return fmt.Errorf("user %s is not an admin of conversation %s", addedByUserID, conversationID)
	}

	isPrivate, err := r.IsUserPrivate(ctx, userID)
	if err != nil {
		return err
	}

	if isPrivate {
		areFriends, err := r.AreUsersFriends(ctx, addedByUserID, userID)
		if err != nil {
			return err
		}
		if !areFriends {
			return nil
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
func (r *Repository) IsUserAdmin(ctx context.Context, conversationID, userID uuid.UUID) (bool, error) {
	var isAdmin bool
	query := `
		SELECT is_admin
		FROM conversation_participants
		WHERE conversation_id = $1 AND user_id = $2
	`
	err := r.db.QueryRowContext(ctx, query, conversationID, userID).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return isAdmin, nil
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
func (r *Repository) GetParticipants(ctx context.Context, conversationID uuid.UUID) ([]domain.User, error) {
	query := `
		SELECT u.id, u.username, u.avatar_url
		FROM users_cache u
		JOIN conversation_participants cp ON cp.user_id = u.id
		WHERE cp.conversation_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []domain.User
	for rows.Next() {
		var user domain.User
		var avatar sql.NullString
		if err := rows.Scan(&user.ID, &user.Username, &avatar); err != nil {
			return nil, err
		}
		if avatar.Valid {
			user.Avatar = avatar.String
		}
		participants = append(participants, user)
	}

	return participants, nil
}
func (r *Repository) PromoteToAdmin(ctx context.Context, conversationID, targetUserID, currentUserID uuid.UUID) error {
	isAdmin, err := r.IsUserAdmin(ctx, conversationID, currentUserID)
	if err != nil {
		return fmt.Errorf("failed to check admin rights: %w", err)
	}
	if !isAdmin {
		return fmt.Errorf("unauthorized: only admins can promote others")
	}

	// Şimdi hedef kullanıcıyı admin yap
	query := `
		UPDATE conversation_participants
		SET is_admin = true
		WHERE conversation_id = $1 AND user_id = $2
	`
	res, err := r.db.ExecContext(ctx, query, conversationID, targetUserID)
	if err != nil {
		return fmt.Errorf("failed to promote user to admin: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found in the conversation")
	}

	return nil
}
func (r *Repository) DemoteFromAdmin(ctx context.Context, conversationID, targetUserID, currentUserID uuid.UUID) error {
	// Önce currentUserID gerçekten admin mi diye kontrol et
	isAdmin, err := r.IsUserAdmin(ctx, conversationID, currentUserID)
	if err != nil {
		return fmt.Errorf("failed to check admin rights: %w", err)
	}
	if !isAdmin {
		return fmt.Errorf("unauthorized: only admins can demote others")
	}

	// Kendini düşürmeyi engelle (opsiyonel ama önerilir)
	if currentUserID == targetUserID {
		return fmt.Errorf("admins cannot demote themselves")
	}

	// Hedef kullanıcıyı adminlikten düşür
	query := `
		UPDATE conversation_participants
		SET is_admin = false
		WHERE conversation_id = $1 AND user_id = $2
	`
	res, err := r.db.ExecContext(ctx, query, conversationID, targetUserID)
	if err != nil {
		return fmt.Errorf("failed to demote user from admin: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found in the conversation")
	}

	return nil
}

func (r *Repository) RemoveParticipant(ctx context.Context, conversationID, userID, addedByUserID uuid.UUID) error {

	isAdderAdmin, err := r.IsUserAdmin(ctx, conversationID, addedByUserID)
	if err != nil {
		return fmt.Errorf("failed to check admin rights: %w", err)
	}
	if !isAdderAdmin {
		return fmt.Errorf("user %s is not an admin of conversation %s", addedByUserID, conversationID)
	}

	query := `
	DELETE FROM conversation_participants
	WHERE conversation_id = $1 AND user_id = $2
`
	_, err = r.db.ExecContext(ctx, query, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}
	return nil
}
