package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"socialmedia/follow/domain"

	"github.com/google/uuid"
)

func (r *Repository) GetBlockedUsers(ctx context.Context, currentUserID uuid.UUID) ([]*domain.BlockedUser, error) {
	query := `
		SELECT 
			b.id AS block_id,
			u.id AS user_id,
			u.username,
			u.avatar_url
		FROM blocks b
		JOIN users_cache u ON b.blocked_id = u.id
		WHERE b.blocker_id = $1
		ORDER BY b.blocked_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get blocked users: %w", err)
	}
	defer rows.Close()

	var blockedUsers []*domain.BlockedUser
	for rows.Next() {
		var blockedUser domain.BlockedUser
		var avatarURL sql.NullString

		if err := rows.Scan(
			&blockedUser.ID,
			&blockedUser.UserID,
			&blockedUser.Username,
			&avatarURL,
		); err != nil {
			return nil, fmt.Errorf("failed to scan blocked user row: %w", err)
		}

		if avatarURL.Valid {
			blockedUser.AvatarURL = avatarURL.String
		}

		blockedUsers = append(blockedUsers, &blockedUser)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating blocked user rows: %w", err)
	}

	return blockedUsers, nil
}
