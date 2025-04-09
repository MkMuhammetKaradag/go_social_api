package auth

import (
	"context"
	"socialmedia/auth/app/auth/usecase"
	"socialmedia/auth/domain"
)

type ActivateAuthRequest struct {
	ActivationToken string `json:"activationToken" binding:"required"`
	ActivationCode  string `json:"activationCode" binding:"required,min=4,max=4"`
}

type ActivateAuthResponse struct {
	Auth *domain.AuthResponse `json:"auth"`
}
type ActivateAuthHandler struct {
	usecase usecase.ActivateUseCase
}

func NewActivateAuthHandler(usecase usecase.ActivateUseCase) *ActivateAuthHandler {
	return &ActivateAuthHandler{
		usecase: usecase,
	}
}

func (h *ActivateAuthHandler) Handle(ctx context.Context, req *ActivateAuthRequest) (*ActivateAuthResponse, error) {
	auth, err := h.usecase.Execute(ctx, req.ActivationToken, req.ActivationCode)
	if err != nil {
		return nil, err
	}

	return &ActivateAuthResponse{Auth: auth}, nil
}
