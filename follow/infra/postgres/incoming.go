package postgres

import (
	"context"
	"fmt"
	"socialmedia/follow/domain"

	"github.com/google/uuid"
)

func (r *Repository) IncomingRequests(ctx context.Context, currentUserID uuid.UUID) ([]*domain.User, error) {
	const query = `
	SELECT 
        u.id,
        u.username,
        fr.requested_at
    FROM follow_requests fr
    JOIN users_cache u ON fr.requester_id = u.id
    WHERE fr.target_id = $1 AND fr.status = 'pending'
	`
	fmt.Println(currentUserID)
	rows, err := r.db.QueryContext(ctx, query, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch incoming requests: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.RequestedAt,
			// &user.AvatarURL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}
