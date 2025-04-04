package auth

import (
	"context"
	"socialmedia/auth/domain"
)

type Repository interface {
	SignUp(ctx context.Context, auth *domain.Auth) error
}
