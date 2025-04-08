package auth

import (
	"context"
	"socialmedia/auth/domain"
)

type Repository interface {
	SignUp(ctx context.Context, auth *domain.Auth) error
	SignIn(ctx context.Context, identifier, password string) (*domain.Auth, error)
	Activate(ctx context.Context, userEmail string, activationCode string) (*domain.Auth, error)
}
