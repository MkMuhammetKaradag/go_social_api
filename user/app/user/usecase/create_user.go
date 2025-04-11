package usecase

import (
	"context"
	"fmt"
)

type createUserUseCase struct {
	repository Repository
}

func NewCreateUserUseCase(repository Repository) CreateUserUseCase {
	return &createUserUseCase{
		repository: repository,
	}
}

func (u *createUserUseCase) Execute(ctx context.Context, userID, userName, email string) error {
	fmt.Println("user created geldi")
	err := u.repository.CreateUser(ctx, userID, userName, email)
	if err != nil {
		return  err

	}

	return nil
}
