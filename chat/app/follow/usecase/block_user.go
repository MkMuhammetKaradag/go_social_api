package usecase

import (
	"context"

	"github.com/google/uuid"
)

type blockUserUseCase struct {
	repository Repository
}

func NewBlockUserUseCase(repository Repository) BlockUserUseCase {
	return &blockUserUseCase{
		repository: repository,
	}
}

func (u *blockUserUseCase) Execute(ctx context.Context, blockerID, blockedID uuid.UUID) error {

	err := u.repository.BlockUser(ctx, blockerID, blockedID)
	if err != nil {
		return err
	}

	return nil

}
