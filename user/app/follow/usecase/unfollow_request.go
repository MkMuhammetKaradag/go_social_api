package usecase

import (
	"context"

	"github.com/google/uuid"
)

type unfollowRequestUseCase struct {
	repository Repository
}

func NewUnFollowRequestUseCase(repository Repository) UnFollowRequestUseCase {
	return &unfollowRequestUseCase{
		repository: repository,
	}
}

func (u *unfollowRequestUseCase) Execute(ctx context.Context, followerID, followingID uuid.UUID, status string) error {

	err := u.repository.CreateFollow(ctx, followerID, followingID, status)
	if err != nil {
		return err
	}

	return nil

}
