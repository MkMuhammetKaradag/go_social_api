package postgres

import (
	"context"
	"fmt"
	"socialmedia/chat/domain"

	"github.com/google/uuid"
)

func (r *Repository) MarkMessagesAsRead(ctx context.Context, messageIDs []uuid.UUID, userID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	getMessageInfoStmt := `
		SELECT sender_id, conversation_id FROM messages WHERE id = $1
	`

	checkParticipantStmt := `
		SELECT EXISTS (
			SELECT 1 FROM conversation_participants
			WHERE conversation_id = $1 AND user_id = $2
		)
	`

	insertStmt := `
		INSERT INTO message_reads (message_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`

	for _, msgID := range messageIDs {
		var senderID, conversationID uuid.UUID

		err := tx.QueryRowContext(ctx, getMessageInfoStmt, msgID).Scan(&senderID, &conversationID)
		if err != nil {
			return fmt.Errorf("failed to fetch message info for message %s: %w", msgID, err)
		}

		// Kendi mesajını okundu olarak işaretleme
		if senderID == userID {
			continue
		}

		// Kullanıcı bu sohbetin katılımcısı mı?
		var isParticipant bool
		err = tx.QueryRowContext(ctx, checkParticipantStmt, conversationID, userID).Scan(&isParticipant)
		if err != nil {
			return fmt.Errorf("failed to check participant for message %s: %w", msgID, err)
		}
		if !isParticipant {
			continue // yetkisi olmayan mesajı geç
		}

		// Okundu olarak işaretle
		if _, err := tx.ExecContext(ctx, insertStmt, msgID, userID); err != nil {
			return fmt.Errorf("failed to mark message %s as read: %w", msgID, err)
		}
	}

	return tx.Commit()
}

func (r *Repository) MarkConversationMessagesAsRead(ctx context.Context, conversationID, userID uuid.UUID) error {
	// Kullanıcının sohbetin katılımcısı olup olmadığını kontrol et
	checkParticipantStmt := `
		SELECT EXISTS (
			SELECT 1 FROM conversation_participants
			WHERE conversation_id = $1 AND user_id = $2
		)
	`

	var isParticipant bool
	err := r.db.QueryRowContext(ctx, checkParticipantStmt, conversationID, userID).Scan(&isParticipant)
	if err != nil {
		return fmt.Errorf("failed to check if user is a participant: %w", err)
	}
	if !isParticipant {
		return fmt.Errorf("user is not a participant in this conversation")
	}

	// Mesajları okundu olarak işaretleme
	query := `
		INSERT INTO message_reads (message_id, user_id)
		SELECT m.id, $2
		FROM messages m
		LEFT JOIN message_reads mr
		  ON m.id = mr.message_id AND mr.user_id = $2
		WHERE m.conversation_id = $1
		  AND mr.message_id IS NULL
		  AND m.sender_id != $2
	`

	_, err = r.db.ExecContext(ctx, query, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark conversation messages as read: %w", err)
	}

	return nil
}
func (r *Repository) GetMessagesForConversation(ctx context.Context, conversationID, userID uuid.UUID, skip, limit int64) ([]domain.Message, error) {
	query := `
	WITH limited_messages AS (
  SELECT m.*
  FROM messages m
  WHERE m.conversation_id = $1
    AND m.deleted_at IS NULL
    AND NOT EXISTS (
      SELECT 1 FROM blocks_cache b
      WHERE b.blocker_id = m.sender_id AND b.blocked_id = $2
    )
    AND EXISTS (
      SELECT 1 FROM conversation_participants cp
      WHERE cp.conversation_id = $1 AND cp.user_id = $2
    )
  ORDER BY m.created_at DESC
  OFFSET $3
  LIMIT $4
)
SELECT
  m.id,
  m.conversation_id,
  m.sender_id,
  u.username,
  u.avatar_url,
  m.content,
  m.created_at,
  m.is_edited,
  m.deleted_at,
  r.read_at,
  a.id AS attachment_id,
  a.file_url,
  a.file_type
FROM limited_messages m
LEFT JOIN users_cache u ON m.sender_id = u.id
LEFT JOIN message_reads r ON m.id = r.message_id AND r.user_id = $2
LEFT JOIN attachments a ON m.id = a.message_id
ORDER BY m.created_at DESC, a.id
	`

	rows, err := r.db.QueryContext(ctx, query, conversationID, userID, skip, limit)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	msgMap := make(map[uuid.UUID]*domain.Message)

	for rows.Next() {
		var (
			msgID uuid.UUID
			msg   domain.Message
			att   domain.Attachment
			// attID sql.NullString
		)

		err := rows.Scan(
			&msgID,
			&msg.ConversationID,
			&msg.UserID,
			&msg.SenderUsername,
			&msg.SenderAvatar,
			&msg.Content,
			&msg.CreatedAt,
			&msg.IsEdited,
			&msg.DeletedAt,
			&msg.ReadAt,
			&att.ID,
			&att.FileURL,
			&att.FileType,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		if existingMsg, exists := msgMap[msgID]; exists {
			if att.ID != uuid.Nil {
				existingMsg.Attachments = append(existingMsg.Attachments, att)
			}
		} else {
			msg.ID = msgID
			if att.ID != uuid.Nil {
				msg.Attachments = []domain.Attachment{att}
			}
			msgMap[msgID] = &msg
		}
	}

	var messages []domain.Message
	for _, msg := range msgMap {
		messages = append(messages, *msg)
	}

	return messages, nil
}
