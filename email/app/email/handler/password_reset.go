package handler

import (
	"socialmedia/email/app/email/usecase"
	"socialmedia/email/internal/domain"
	"socialmedia/shared/messaging"
)

type PasswordResetHandler struct {
	passwordResetUsecase usecase.PasswordResetUseCase
}

func NewPasswordResetEmailHandler(passwordResetUsecase usecase.PasswordResetUseCase) *PasswordResetHandler {
	return &PasswordResetHandler{
		passwordResetUsecase: passwordResetUsecase,
	}
}

func (h *PasswordResetHandler) HandleEmail(msg messaging.Message) error {

	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return domain.ErrInvalidMessageFormat
	}

	email, ok := data["email"].(string)
	if !ok {
		return domain.ErrMissingEmail
	}

	resetLink, ok := data["reset_link"].(string)
	if !ok {
		return domain.ErrMissingResetLink
	}

	userName, _ := data["userName"].(string)
	templateName, ok := data["template_name"].(string)
	if !ok {
		return domain.ErrMissingTemplateName
	}

	return h.passwordResetUsecase.SendPasswordResetEmail(email, resetLink, userName, templateName)
}
