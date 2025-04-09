package usecase

import (
	"context"
	"socialmedia/auth/domain"

	"github.com/gofiber/fiber/v2"
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
