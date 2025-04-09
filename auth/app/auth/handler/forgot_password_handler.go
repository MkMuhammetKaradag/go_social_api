package auth

import (
	"context"
	"socialmedia/auth/app/auth/usecase"
)

type ForgotPasswordAuthHandler struct {
	usecase usecase.ForgotPasswordUseCase
}
type ForgotPasswordAuthRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordAuthResponse struct {
	Message string `json:"message"`
}

func NewForgotPasswordAuthHandler(usecase usecase.ForgotPasswordUseCase) *ForgotPasswordAuthHandler {
	return &ForgotPasswordAuthHandler{usecase: usecase}
}

func (h *ForgotPasswordAuthHandler) Handle(ctx context.Context, req *ForgotPasswordAuthRequest) (*ForgotPasswordAuthResponse, error) {
	err := h.usecase.Execute(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return &ForgotPasswordAuthResponse{Message: "Reset Password token send email"}, nil
}
