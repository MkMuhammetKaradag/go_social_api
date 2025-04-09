package usecase

import (
	"context"
	"errors"
	"fmt"
	"socialmedia/auth/domain"
)

type activateUseCase struct {
	repository Repository
	jwtHelper  JwtHelper
}

func NewActivateUseCase(repository Repository, jwtHelper JwtHelper) ActivateUseCase {
	return &activateUseCase{
		repository: repository,
		jwtHelper:  jwtHelper,
	}
}

func (u *activateUseCase) Execute(ctx context.Context, activationToken, activationCode string) (*domain.AuthResponse, error) {
	claims, err := u.jwtHelper.VerifyToken(activationToken)
	if err != nil {
		return nil, fmt.Errorf("error verifying token: %w", err)
	}
	// code := claims["activationCode"].(string)
	// if code != activationCode {
	// 	return nil, errors.New("activation code mismatch")
	// }
	userEmail, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid user data in token")
	}

	auth, err := u.repository.Activate(ctx, userEmail, activationCode)
	if err != nil {
		return nil, err
	}

	response := &domain.AuthResponse{
		ID:       auth.ID,
		Username: auth.Username,
		Email:    auth.Email,
	}
	return response, nil
}
