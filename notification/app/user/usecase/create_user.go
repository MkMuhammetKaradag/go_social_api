package usecase

import (
	"context"

	"github.com/google/uuid"
)

type createUserUseCase struct {
	repository      Repository
	repositoryMongo RepositoryMongo
}

func NewCreateUserUseCase(repository Repository, repositoryMongo RepositoryMongo) CreateUserUseCase {
	return &createUserUseCase{
		repository:      repository,
		repositoryMongo: repositoryMongo,
	}
}

func (u *createUserUseCase) Execute(ctx context.Context, userID uuid.UUID, userName string) error {

	err := u.repository.CreateUser(ctx, userID, userName)
	if err != nil {
		return err

	}
	err = u.repositoryMongo.CreateUser(ctx, userID, userName)
	if err != nil {
		return err

	}

	return nil
}
