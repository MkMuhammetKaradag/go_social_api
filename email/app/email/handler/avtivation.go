package handler

import (
	"socialmedia/email/app/email/usecase"
	"socialmedia/email/internal/domain"
	"socialmedia/shared/messaging"
)

type ActivationHandler struct {
	activationUsecase usecase.ActivationUseCase
}

func NewActivationEmailHandler(activationUsecase usecase.ActivationUseCase) *ActivationHandler {
	return &ActivationHandler{
		activationUsecase: activationUsecase,
	}
}

func (h *ActivationHandler) HandleEmail(msg messaging.Message) error {

	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return domain.ErrInvalidMessageFormat
	}

	email, ok := data["email"].(string)
	if !ok {
		return domain.ErrMissingEmail
	}

	code, ok := data["activation_code"].(string)
	if !ok {
		return domain.ErrMissingActivationCode
	}

	userName, _ := data["userName"].(string)
	templateName, ok := data["template_name"].(string)
	if !ok {
		return domain.ErrMissingTemplateName
	}

	return h.activationUsecase.SendActivationEmail(email, code, userName, templateName)
}
