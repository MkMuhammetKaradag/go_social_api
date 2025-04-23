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

func (u *unfollowRequestUseCase) Execute(ctx context.Context, followerID, followingID uuid.UUID) error {

	err := u.repository.DeleteFollow(ctx, followerID, followingID)
	if err != nil {
		return err
	}

	return nil

}
