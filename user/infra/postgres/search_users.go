package postgres

import (
	"context"
	"fmt"
	"socialmedia/user/domain"

	"github.com/google/uuid"
)

func (r *Repository) SearchUsers(ctx context.Context, currentUserID uuid.UUID, searchTerm string, page, limit int) ([]*domain.UserSearchResult, error) {
	offset := (page - 1) * limit

	const query = `
	SELECT
		u.id,
		u.username,
		CASE
			WHEN u.is_private AND NOT COALESCE(f.status = 'following', false) AND u.id != $1 THEN NULL
			ELSE u.avatar_url
		END as avatar_url
	FROM users u
	LEFT JOIN follows_cache f
		ON f.follower_id = $1 AND f.following_id = u.id
	WHERE
		u.username ILIKE $2
		AND NOT EXISTS (
			SELECT 1 FROM blocks_cache
			WHERE (blocker_id = $1 AND blocked_id = u.id)
			   OR (blocker_id = u.id AND blocked_id = $1)
		)
	ORDER BY 
		CASE WHEN u.username ILIKE $3 THEN 0 ELSE 1 END, 
		CASE WHEN u.username ILIKE $4 THEN 0 ELSE 1 END,
		u.followers_count DESC                           
	LIMIT $5 OFFSET $6
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		currentUserID,
		"%"+searchTerm+"%",
		searchTerm,
		searchTerm+"%",
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users []*domain.UserSearchResult
	for rows.Next() {
		var user domain.UserSearchResult
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.AvatarURL,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}
