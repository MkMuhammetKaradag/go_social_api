package usecase

import (
	"context"
	"fmt"
	"socialmedia/shared/middlewares"
	"socialmedia/user/domain"

	"github.com/gofiber/fiber/v2"
)

type updateUserUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
}

func NewUpdateUserUseCase(sessionRepo RedisRepository, repository Repository) UpdateUserUseCase {
	return &updateUserUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
	}
}

func (u *updateUserUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, updateuser domain.UserUpdate) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return  fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}
	userID := userData["id"]
	err := u.repository.UpdateUser(ctx, userID ,updateuser)
	if err != nil {
		return err

	}

	return nil
}
