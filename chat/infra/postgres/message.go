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

	checkSenderStmt := `
		SELECT sender_id FROM messages WHERE id = $1
	`

	insertStmt := `
		INSERT INTO message_reads (message_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`

	for _, msgID := range messageIDs {
		var senderID uuid.UUID
		err := tx.QueryRowContext(ctx, checkSenderStmt, msgID).Scan(&senderID)
		if err != nil {
			return fmt.Errorf("failed to fetch sender for message %s: %w", msgID, err)
		}

		if senderID == userID {
			continue
		}

		if _, err := tx.ExecContext(ctx, insertStmt, msgID, userID); err != nil {
			return fmt.Errorf("failed to mark message %s as read: %w", msgID, err)
		}
	}

	return tx.Commit()
}

