package auth

import (
	"context"
	"socialmedia/auth/domain"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type SignInAuthRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required,min=8"`
}

type SignInAuthResponse struct {
	Auth *domain.Auth `json:"auth"`
}

type SignInAuthHandler struct {
	repository  Repository
	sessionRepo RedisRepository
}

func NewSignInAuthHandler(repository Repository, sessionRepo RedisRepository) *SignInAuthHandler {
	return &SignInAuthHandler{
		repository:  repository,
		sessionRepo: sessionRepo,
	}
}

func (h *SignInAuthHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *SignInAuthRequest) (*SignInAuthResponse, error) {
	auth, err := h.repository.SignIn(ctx, req.Identifier, req.Password)
	if err != nil {
		return nil, err
	}
	sessionKey := "session:" + strconv.Itoa(int(auth.ID))
	userData := map[string]string{
		"id":       strconv.Itoa(int(auth.ID)),
		"email":    auth.Email,
		"username": auth.Username,
	}

	// ctx := context.Background()
	if err := h.sessionRepo.SetSession(ctx, sessionKey, userData, 24*time.Hour); err != nil {
		return nil, err
	}
	fbrCtx.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionKey,
		Path:     "/",
		MaxAge:   60 * 60 * 24, // 1 gün
		HTTPOnly: true,
		Secure:   false, // HTTPS kullanıyorsan true yap
		SameSite: "Lax",
	})
	return &SignInAuthResponse{Auth: auth}, nil
}
