package postgres

import (
	"context"
	"fmt"

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
