package usecase

import (
	"context"
	"errors"
	"fmt"
	"socialmedia/auth/domain"
	"socialmedia/shared/messaging"
)

type activateUseCase struct {
	repository Repository
	jwtHelper  JwtHelper
	rabbitMQ   RabbitMQ
}

func NewActivateUseCase(repository Repository, jwtHelper JwtHelper, rabbitMQ RabbitMQ) ActivateUseCase {
	return &activateUseCase{
		repository: repository,
		jwtHelper:  jwtHelper,
		rabbitMQ:   rabbitMQ,
	}
}

func (u *activateUseCase) Execute(ctx context.Context, activationToken, activationCode string) (*domain.AuthResponse, error) {
	claims, err := u.jwtHelper.VerifyToken(activationToken)
	if err != nil {
		return nil, fmt.Errorf("error verifying token: %w", err)
	}
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

	userCreatedMessage := messaging.Message{
		Type:       messaging.UserTypes.UserCreated,
		ToService:  messaging.UserService,
		RetryCount: 4,
		Data: map[string]interface{}{
			"id":       response.ID,
			"email":    response.Email,
			"username": response.Username,
		},
	}

	if err := u.rabbitMQ.PublishMessage(context.Background(), userCreatedMessage); err != nil {
		// log.Printf("User creation message could not be sent: %v", err)
		return nil, err
	}

	return response, nil
}
