package auth

import (
	"context"
	"socialmedia/auth/app/auth/usecase"
)

type ResetPasswordAuthRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type ResetPasswordAuthResponse struct {
	Message string `json:"message" `
}
type ResetPasswordAuthHandler struct {
	usecase usecase.ResetPasswordUseCase
}

func NewResetPasswordAuthHandler(usecase usecase.ResetPasswordUseCase) *ResetPasswordAuthHandler {
	return &ResetPasswordAuthHandler{
		usecase: usecase,
	}
}

func (h *ResetPasswordAuthHandler) Handle(ctx context.Context, req *ResetPasswordAuthRequest) (*ResetPasswordAuthResponse, error) {
	err := h.usecase.Execute(ctx, req.Token, req.Password)
	if err != nil {
		return nil, err
	}

	return &ResetPasswordAuthResponse{Message: "Password changed successfully"}, nil
}
