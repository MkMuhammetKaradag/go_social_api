package usecase

import (
	"context"
	"fmt"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type updateBannerUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
}

func NewUpdateBannerUseCase(sessionRepo RedisRepository, repository Repository) UpdateBannerUseCase {
	return &updateBannerUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
	}
}

func (u *updateBannerUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, avatarURL string) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}
	// userID := userData["id"]
	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}
	err = u.repository.UpdateBanner(ctx, currrentUserID, avatarURL)
	if err != nil {
		return err

	}

	return nil
}
