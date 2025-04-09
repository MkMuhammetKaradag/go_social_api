package usecase

import (
	"context"
	"fmt"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
)

type logoutUseCase struct {
	sessionRepo RedisRepository
}

func NewLogoutUseCase(sessionRepo RedisRepository) LogoutUseCase {
	return &logoutUseCase{
		sessionRepo: sessionRepo,
	}
}

func (u *logoutUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context) error {
	cookieSessionId := fbrCtx.Cookies("session_id")
	if cookieSessionId == "" {
		return fmt.Errorf("user not logout")
	}

	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}

	userID := userData["id"]
	if userID == "" {
		return fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}

	if err := u.sessionRepo.DeleteSession(ctx, cookieSessionId, userID); err != nil {
		return err
	}

	fbrCtx.ClearCookie("session_id")

	return nil
}
