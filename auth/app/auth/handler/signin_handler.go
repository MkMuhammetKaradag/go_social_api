package auth

import (
	"context"
	"socialmedia/auth/app/auth/usecase"
	"socialmedia/auth/domain"
	"github.com/gofiber/fiber/v2"
)

type SignInAuthRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required,min=8"`
}

type SignInAuthResponse struct {
	Auth *domain.AuthResponse `json:"auth"`
}
type SignInAuthHandler struct {
	usecase usecase.SignInUseCase
}

func NewSignInAuthHandler(usecase usecase.SignInUseCase) *SignInAuthHandler {
	return &SignInAuthHandler{
		usecase: usecase,
	}
}

func (h *SignInAuthHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *SignInAuthRequest) (*SignInAuthResponse, error) {
	auth, err := h.usecase.Execute(fbrCtx, ctx, req.Identifier, req.Password)
	if err != nil {
		return nil, err
	}

	return &SignInAuthResponse{Auth: auth}, nil
}
