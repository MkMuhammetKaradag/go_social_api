package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) DeleteAllMessagesFromConversation(ctx context.Context, conversationID, currentUserID uuid.UUID) error {

	isAdmin, err := r.IsUserAdmin(ctx, conversationID, currentUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user is not a participant in this conversation")
		}
		return fmt.Errorf("failed to check admin status: %w", err)
	}
	if !isAdmin {
		return fmt.Errorf("only admins can delete all messages in the conversation")
	}
	// Kullanıcı admin, mesajları silebilir
	deleteQuery := `
		DELETE FROM messages
		WHERE conversation_id = $1
	`

	_, err = r.db.ExecContext(ctx, deleteQuery, conversationID)
	if err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}

	return nil
}
