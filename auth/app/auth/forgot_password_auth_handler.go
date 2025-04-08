package auth

import (
	"context"
	"socialmedia/auth/domain"
	"socialmedia/shared/messaging"
	"time"

	"github.com/google/uuid"
)

type ForgotPasswordAuthRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordAuthResponse struct {
	Message string `json:"message"`
}

type ForgotPasswordAuthHandler struct {
	repository Repository
	rabbitMQ   RabbitMQ
}

func NewForgotPasswordAuthHandler(repository Repository, rabbitMQ RabbitMQ) *ForgotPasswordAuthHandler {
	return &ForgotPasswordAuthHandler{
		repository: repository,
		rabbitMQ:   rabbitMQ,
	}
}

func (h *ForgotPasswordAuthHandler) Handle(ctx context.Context, req *ForgotPasswordAuthRequest) (*ForgotPasswordAuthResponse, error) {

	token, resetPassword := h.generateForgotPasswordLink(req.Email)
	username, err := h.repository.RequestForgotPassword(ctx, resetPassword)
	if err != nil {
		return nil, err
	}

	emailMessage := messaging.Message{
		Type:      messaging.EmailTypes.ForgotPassword,
		ToService: messaging.EmailService,
		Data: map[string]interface{}{
			"email":           req.Email,
			"activation_code": token,
			"template_name":   "activation_email.html",
			"userName":        *username,
		},
	}

	if err := h.rabbitMQ.PublishMessage(context.Background(), emailMessage); err != nil {

		return nil, err
	}

	return &ForgotPasswordAuthResponse{Message: "Reset Password token send email"}, nil
}

func (h *ForgotPasswordAuthHandler) generateForgotPasswordLink(email string) (string, *domain.ForgotPassword) {
	resetToken := uuid.NewString()
	expiresAt := time.Now().Add(1 * time.Hour)

	passwordReset := &domain.ForgotPassword{
		Email:     email,
		Token:     resetToken,
		ExpiresAt: expiresAt,
	}

	return resetToken, passwordReset
}
