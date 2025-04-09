package usecase

import (
	"context"
	"socialmedia/auth/domain"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type signInUseCase struct {
	repository  Repository
	sessionRepo RedisRepository
}

func NewSignInUseCase(repository Repository, sessionRepo RedisRepository) SignInUseCase {
	return &signInUseCase{
		repository:  repository,
		sessionRepo: sessionRepo,
	}
}

func (u *signInUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, identifier, password string) (*domain.AuthResponse, error) {
	auth, err := u.repository.SignIn(ctx, identifier, password)
	if err != nil {
		return nil, err
	}
	sessionKey := strconv.Itoa(int(auth.ID))
	sessionID := uuid.New().String()
	device := fbrCtx.Get("User-Agent")
	ip := fbrCtx.IP()

	// fmt.Println(device, ip)
	userData := map[string]string{
		"id":       strconv.Itoa(int(auth.ID)),
		"email":    auth.Email,
		"device":   device,
		"ip":       ip,
		"username": auth.Username,
	}
	if err := u.sessionRepo.SetSession(ctx, sessionID, sessionKey, userData, 24*time.Hour); err != nil {
		return nil, err
	}
	fbrCtx.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})
	response := &domain.AuthResponse{
		ID:       auth.ID,
		Username: auth.Username,
		Email:    auth.Email,
	}
	return response, nil
}
