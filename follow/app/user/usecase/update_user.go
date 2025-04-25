package usecase

import (
	"context"

	"github.com/google/uuid"
)

type updateUserUseCase struct {
	repository Repository
}

func NewUpdateUserUseCase(repository Repository) UpdateUserUseCase {
	return &updateUserUseCase{
		repository: repository,
	}
}

func (u *updateUserUseCase) Execute(ctx context.Context, userID uuid.UUID, userName, avatarURL *string, isPrivate *bool) error {

	err := u.repository.UpdateUser(ctx, userID, userName, avatarURL, isPrivate)
	if err != nil {
		return err

	}

	return nil
}
