package usecase

import (
	"context"
	"strconv"
)

type resetPasswordUseCase struct {
	repository  Repository
	sessionRepo RedisRepository
}

func NewResetPasswordUseCase(repo Repository, sessionRepo RedisRepository) ResetPasswordUseCase {
	return &resetPasswordUseCase{
		repository:  repo,
		sessionRepo: sessionRepo,
	}
}

func (u *resetPasswordUseCase) Execute(ctx context.Context, token, password string) error {
	userID, err := u.repository.ResetPassword(ctx, token, password)
	if err != nil {
		return err
	}
	strUserID := strconv.Itoa(*userID)
	if err := u.sessionRepo.DeleteAllUserSessions(ctx, strUserID); err != nil {
		return err
	}
	return nil

}
