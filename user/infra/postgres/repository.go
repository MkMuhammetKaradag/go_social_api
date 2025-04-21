package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"socialmedia/user/domain"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
	ErrUserNotFound = errors.New("user not found")
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

func (r *Repository) GetUserProfile(ctx context.Context, identifier string) (*domain.User, error) {
	var query string
	var row *sql.Row

	// Eğer UUID mi kontrol et
	if _, err := uuid.Parse(identifier); err == nil {
		// UUID ise id ile çek
		query = `SELECT id, username, email, bio, avatar_url, is_private, created_at, updated_at FROM users WHERE id = $1`
		row = r.db.QueryRowContext(ctx, query, identifier)
	} else {
		// Değilse username/email ile çek
		query = `SELECT id, username, email, bio, avatar_url, is_private, created_at, updated_at FROM users WHERE username = $1 OR email = $1`
		row = r.db.QueryRowContext(ctx, query, identifier)
	}

	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Bio,
		&user.AvatarURL,
		&user.IsPrivate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("query error: %w", err)
	}
	// fmt.Println("repository içi", user)

	return &user, nil
}

func (r *Repository) CreateUser(ctx context.Context, id, username, email string) error {

	query := `
		INSERT INTO users (id, username, email)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
        SET username = EXCLUDED.username,
            email = EXCLUDED.email
	`
	_, err := r.db.ExecContext(ctx, query, id, username, email)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *Repository) GetUser(ctx context.Context, currrentUserID, targetUserID uuid.UUID) (*domain.User, error) {
	const query = `
SELECT 
	u.id,
	u.username,
	u.avatar_url,
	u.banner_url,
	u.is_private,
	COALESCE(f.status = 'following', false) AS is_following,
	EXISTS (
		SELECT 1 FROM blocks_cache 
		WHERE (blocker_id = $1 AND blocked_id = $2)
		   OR (blocker_id = $2 AND blocked_id = $1)
	) AS is_blocked,
	u.id = $2 AS is_self
FROM users u
LEFT JOIN follows_cache f 
	ON f.follower_id = $2 AND f.following_id = u.id
WHERE u.id = $1

`

	row := r.db.QueryRowContext(ctx, query, targetUserID, currrentUserID)

	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.AvatarURL,
		&user.BannerURL,
		&user.IsPrivate,
		&user.IsFollowing,
		&user.IsBlocked,
		&user.Self,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
func (r *Repository) UpdateUser(ctx context.Context, userID string, update domain.UserUpdate) error {
	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if update.Bio != nil {
		setClauses = append(setClauses, fmt.Sprintf("bio = $%d", argPos))
		args = append(args, *update.Bio)
		argPos++
	}
	if update.AvatarURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("avatar_url = $%d", argPos))
		args = append(args, *update.AvatarURL)
		argPos++
	}

	if update.BannerURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("banner_url = $%d", argPos))
		args = append(args, *update.BannerURL)
		argPos++
	}
	if update.Location != nil {
		setClauses = append(setClauses, fmt.Sprintf("location = $%d", argPos))
		args = append(args, *update.Location)
		argPos++
	}
	if update.Website != nil {
		setClauses = append(setClauses, fmt.Sprintf("website = $%d", argPos))
		args = append(args, *update.Website)
		argPos++
	}
	if update.IsPrivate != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_private = $%d", argPos))
		args = append(args, *update.IsPrivate)
		argPos++
	}

	// Eğer hiçbir alan gönderilmemişse, boş update yapma:
	if len(setClauses) == 0 {
		return errors.New("no fields to update")
	}

	// updated_at alanını da güncelleyelim:
	setClauses = append(setClauses, fmt.Sprintf("updated_at = CURRENT_TIMESTAMP"))

	// WHERE id = $N
	query := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE id = $%d
	`,
		strings.Join(setClauses, ", "),
		argPos,
	)

	args = append(args, userID)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
