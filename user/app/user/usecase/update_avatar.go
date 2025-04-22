package usecase

import (
	"context"
	"fmt"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type updateAvatarUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
}

func NewUpdateAvatarUseCase(sessionRepo RedisRepository, repository Repository) UpdateAvatarUseCase {
	return &updateAvatarUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
	}
}

func (u *updateAvatarUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, avatarURL string) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}
	// userID := userData["id"]
	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}
	err = u.repository.UpdateAvatar(ctx, currrentUserID, avatarURL)
	if err != nil {
		return err

	}

	return nil
}
