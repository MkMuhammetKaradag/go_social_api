package bootstrap

import (
	"context"
	"socialmedia/auth/domain"
	"socialmedia/auth/internal/initializer"
	"socialmedia/auth/pkg/config"
	"time"

	"github.com/golang-jwt/jwt"
)

type Repository interface {
	SignUp(ctx context.Context, auth *domain.Auth) error
	SignIn(ctx context.Context, identifier, password string) (*domain.Auth, error)
	Activate(ctx context.Context, userEmail string, activationCode string) (*domain.Auth, error)
	RequestForgotPassword(ctx context.Context, resetPassword *domain.ForgotPassword) (string, error)
	ResetPassword(ctx context.Context, token, password string) (*int, error)
}

type RedisRepository interface {
	SetSession(ctx context.Context, key string, userId string, userData map[string]string, expiration time.Duration) error
	GetSession(ctx context.Context, key string) (map[string]string, error)
	DeleteSession(ctx context.Context, key string, userID string) error
	DeleteAllUserSessions(ctx context.Context, userId string) error
}

type JwtHelper interface {
	SignToken(payload jwt.MapClaims, expiration time.Duration) (string, error)
	VerifyToken(tokenStr string) (jwt.MapClaims, error)
}

func InitDatabase(config config.Config) Repository {
	return initializer.InitDatabase(config)
}

func InitRedis(config config.Config) RedisRepository {
	return initializer.InitRedis(config)
}

func InitJwtHelper(config config.Config) JwtHelper {
	return initializer.InitJwtHelper(config)
}
