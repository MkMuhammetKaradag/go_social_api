package usecase

import (
	"context"
	"errors"
	"fmt"
	// "errors"
	// "fmt"
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
	// userCreatedMessage := messaging.Message{
	// 	Type:       messaging.UserTypes.UserCreated,
	// 	ToServices: []messaging.ServiceType{messaging.UserService, messaging.EmailService},
	// 	RetryCount: 0,
	// 	Data: map[string]interface{}{
	// 		"id":       "f2c8f899-a358-4805-a8b0-2134e4d68189",
	// 		"email":    "karadag",
	// 		"username": "karadag",
	// 	},
	// }

	// if err := u.rabbitMQ.PublishMessage(context.Background(), userCreatedMessage); err != nil {
	// 	// log.Printf("User creation message could not be sent: %v", err)
	// 	return nil, err
	// }

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
		ToServices: []messaging.ServiceType{messaging.UserService},
		Data: map[string]interface{}{
			"id":       response.ID,
			"email":    response.Email,
			"username": response.Username,
		},
		Critical: true,
	}

	if err := u.rabbitMQ.PublishMessage(context.Background(), userCreatedMessage); err != nil {
		// log.Printf("User creation message could not be sent: %v", err)
		return nil, err
	}
	// response := &domain.AuthResponse{
	// 	ID:       "auth.ID",
	// 	Username: "asas",
	// 	Email:    "sadsadasd",
	// }
	return response, nil
}
