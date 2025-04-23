package usecase

import (
	"context"

	"github.com/google/uuid"
)

type unBlockUserUseCase struct {
	repository Repository
}

func NewUnBlockUserUseCase(repository Repository) UnBlockUserUseCase {
	return &unBlockUserUseCase{
		repository: repository,
	}
}

func (u *unBlockUserUseCase) Execute(ctx context.Context, blockerID, blockedID uuid.UUID) error {

	err := u.repository.UnblockUser(ctx, blockerID, blockedID)
	if err != nil {
		return err
	}

	return nil

}
