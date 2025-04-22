package postgres

import (
	"context"
	"fmt"
	"socialmedia/follow/domain"

	"github.com/google/uuid"
)

func (r *Repository) OutgoingRequests(ctx context.Context, currentUserID uuid.UUID) ([]*domain.FollowRequestUser, error) {
	const query = `
	SELECT 
	  fr.id,
        u.id,
        u.username,
        fr.requested_at
    FROM follow_requests fr
    JOIN users_cache u ON fr.target_id = u.id
    WHERE fr.requester_id = $1 AND fr.status = 'pending'
	`

	rows, err := r.db.QueryContext(ctx, query, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch outgoing requests: %w", err)
	}
	defer rows.Close()

	var requests []*domain.FollowRequestUser
	for rows.Next() {
		var request domain.FollowRequestUser
		err := rows.Scan(
			&request.ID,
			&request.UserID,
			&request.Username,
			&request.RequestedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		requests = append(requests, &request)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return requests, nil
}
