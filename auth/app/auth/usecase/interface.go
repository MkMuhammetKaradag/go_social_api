package usecase

import (
	"context"
	"socialmedia/auth/domain"
	"time"

	"socialmedia/shared/messaging"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type ForgotPasswordUseCase interface {
	Execute(ctx context.Context, email string) error
}
type ResetPasswordUseCase interface {
	Execute(ctx context.Context, token, password string) error
}
type LogoutUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context) error
}
type SignInUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context, identifier, password string) (*domain.AuthResponse, error)
}
type SignUpUseCase interface {
	Execute(ctx context.Context, req *SignUpRequest) (*string, error)
}
type ActivateUseCase interface {
	Execute(ctx context.Context, activationToken, activationCode string) (*domain.AuthResponse, error)
}

type JwtHelper interface {
	SignToken(payload jwt.MapClaims, expiration time.Duration) (string, error)
	VerifyToken(tokenStr string) (jwt.MapClaims, error)
}

type RabbitMQ interface {
	PublishMessage(ctx context.Context, msg messaging.Message) error
}
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
