package usecase

import (
	"context"
	"fmt"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
)

type profileUseCase struct {
	sessionRepo RedisRepository
}

func NewProfileUseCase(sessionRepo RedisRepository) ProfileUseCase {
	return &profileUseCase{
		sessionRepo: sessionRepo,
	}
}

func (u *profileUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context) error {
	cookieSessionId := fbrCtx.Cookies("session_id")
	if cookieSessionId == "" {
		return fmt.Errorf("user not profile")
	}

	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}

	userID := userData["id"]
	if userID == "" {
		return fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}

	fmt.Println(userID)

	return nil
}
