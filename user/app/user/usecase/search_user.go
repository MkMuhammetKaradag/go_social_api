package usecase

import (
	"context"
	"fmt"
	"socialmedia/shared/middlewares"
	"socialmedia/user/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type searchUserUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
}

func NewSearchUserUseCase(sessionRepo RedisRepository, repository Repository) SearchUserUseCase {
	return &searchUserUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
	}
}

func (u *searchUserUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, identifier string, page, limit int) ([]*domain.UserSearchResult, error) {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return nil, fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}
	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return nil, err

	}

	users, err := u.repository.SearchUsers(ctx, currrentUserID, identifier, page, limit)
	if err != nil {
		return nil, err

	}

	return users, nil
}
