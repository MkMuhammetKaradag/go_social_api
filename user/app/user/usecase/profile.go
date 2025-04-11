package usecase

import (
	"context"
	"fmt"
	"socialmedia/shared/middlewares"
	"socialmedia/user/domain"

	"github.com/gofiber/fiber/v2"
)

type profileUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
}

func NewProfileUseCase(sessionRepo RedisRepository, repository Repository) ProfileUseCase {
	return &profileUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
	}
}

func (u *profileUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context) (*domain.User, error) {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return nil, fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}
	userID := userData["id"]
	if userID == "" {
		return nil, fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}
	user, err := u.repository.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err

	}
	fmt.Println(user)

	return user, nil
}
