package auth

import (
	"context"
	"socialmedia/auth/app/auth/usecase"

	"github.com/gofiber/fiber/v2"
)

type AllLogoutAuthRequest struct {
}

type AllLogoutAuthResponse struct {
	Message string `json:"message"`
}

type AllLogoutAuthHandler struct {
	usecase usecase.LogoutUseCase
}

func NewAllLogoutAuthHandler(usecase usecase.LogoutUseCase) *AllLogoutAuthHandler {
	return &AllLogoutAuthHandler{
		usecase: usecase,
	}
}

func (h *AllLogoutAuthHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *AllLogoutAuthRequest) (*AllLogoutAuthResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx)
	if err != nil {
		return nil, err
	}

	return &AllLogoutAuthResponse{Message: "logout user "}, nil
}
