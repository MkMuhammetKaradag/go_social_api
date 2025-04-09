package auth

import (
	"context"
	"socialmedia/auth/app/auth/usecase"

	"github.com/gofiber/fiber/v2"
)

type LogoutAuthRequest struct {
}

type LogoutAuthResponse struct {
	Message string `json:"message"`
}

type LogoutAuthHandler struct {
	usecase usecase.LogoutUseCase
}

func NewLogoutAuthHandler(usecase usecase.LogoutUseCase) *LogoutAuthHandler {
	return &LogoutAuthHandler{
		usecase: usecase,
	}
}

func (h *LogoutAuthHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *LogoutAuthRequest) (*LogoutAuthResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx)
	if err != nil {
		return nil, err
	}

	return &LogoutAuthResponse{Message: "logout user "}, nil
}
