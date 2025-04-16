package usecase

import (
	"context"
)

type createUserUseCase struct {
	repository Repository
}

func NewCreateUserUseCase(repository Repository) CreateUserUseCase {
	return &createUserUseCase{
		repository: repository,
	}
}

func (u *createUserUseCase) Execute(ctx context.Context, userID, userName string) error {

	err := u.repository.CreateUser(ctx, userID, userName)
	if err != nil {
		return err

	}

	return nil
}
