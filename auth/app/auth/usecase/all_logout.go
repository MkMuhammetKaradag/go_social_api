package usecase

import (
	"context"
	"fmt"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type allLogoutUseCase struct {
	sessionRepo   RedisRepository
	userRedisRepo UserRedisRepository
}

func NewAllLogoutUseCase(sessionRepo RedisRepository, userRedisRepo UserRedisRepository) AllLogoutUseCase {
	return &allLogoutUseCase{
		sessionRepo:   sessionRepo,
		userRedisRepo: userRedisRepo,
	}
}

func (u *allLogoutUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context) error {

	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}

	userIDstr := userData["id"]
	if userIDstr == "" {
		return fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}

	if err := u.sessionRepo.DeleteAllUserSessions(ctx, userIDstr); err != nil {
		return err
	}
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return err
	}
	u.userRedisRepo.PublishUserLogout(ctx, userID)

	fbrCtx.ClearCookie("session_id")

	return nil
}
