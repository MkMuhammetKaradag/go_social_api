package usecase

import (
	"context"
	"socialmedia/auth/domain"
	"socialmedia/shared/messaging"
	"time"

	"github.com/google/uuid"
)

type forgotPasswordUseCase struct {
	repository Repository
	rabbitMQ   RabbitMQ
}

func NewForgotPasswordUseCase(repo Repository, rabbitMQ RabbitMQ) ForgotPasswordUseCase {
	return &forgotPasswordUseCase{
		repository: repo,
		rabbitMQ:   rabbitMQ,
	}
}

func (u *forgotPasswordUseCase) Execute(ctx context.Context, email string) error {
	token := uuid.NewString()
	expiresAt := time.Now().Add(1 * time.Hour)

	forgotPassword := &domain.ForgotPassword{
		Email:     email,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	username, err := u.repository.RequestForgotPassword(ctx, forgotPassword)
	if err != nil {
		return err
	}

	emailMessage := messaging.Message{
		Type:      messaging.EmailTypes.ForgotPassword,
		ToService: messaging.EmailService,
		Data: map[string]interface{}{
			"email":           email,
			"activation_code": token,
			"template_name":   "forgot_password.html",
			"userName":        username,
		},
	}

	return u.rabbitMQ.PublishMessage(ctx, emailMessage)
}
