package usecase

import (
	"context"
	"fmt"
	"socialmedia/shared/middlewares"
	"socialmedia/user/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type getUserUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
}

func NewGetUserUseCase(sessionRepo RedisRepository, repository Repository) GetUserUseCase {
	return &getUserUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
	}
}

func (u *getUserUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, identifier uuid.UUID) (*domain.User, error) {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return nil, fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}
	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return nil, err

	}

	user, err := u.repository.GetUser(ctx, currrentUserID, identifier)
	if err != nil {
		return nil, err

	}

	return user, nil
}
