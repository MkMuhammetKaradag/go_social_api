package user

import (
	"context"
	"socialmedia/user/app/user/usecase"

	"github.com/gofiber/fiber/v2"
)

type ProfileUserHandler struct {
	usecase usecase.ProfileUseCase
}
type ProfileUserRequest struct {
}

type ProfileUserResponse struct {
	Message string `json:"message"`
}

func NewProfileUserHandler(usecase usecase.ProfileUseCase) *ProfileUserHandler {
	return &ProfileUserHandler{usecase: usecase}
}

func (h *ProfileUserHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *ProfileUserRequest) (*ProfileUserResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx)
	if err != nil {
		return nil, err
	}
	return &ProfileUserResponse{Message: "get user profile runnig"}, nil
}
