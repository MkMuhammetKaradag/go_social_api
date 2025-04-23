package usecase

import (
	"context"

	"github.com/google/uuid"
)

type followRequestUseCase struct {
	repository Repository
}

func NewFollowRequestUseCase(repository Repository) FollowRequestUseCase {
	return &followRequestUseCase{
		repository: repository,
	}
}

func (u *followRequestUseCase) Execute(ctx context.Context, followerID, followingID uuid.UUID, status string) error {

	err := u.repository.CreateFollow(ctx, followerID, followingID, status)
	if err != nil {
		return err
	}

	return nil

}
