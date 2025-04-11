package user

import (
	"context"
	"socialmedia/user/app/user/usecase"
	"socialmedia/user/domain"

	"github.com/gofiber/fiber/v2"
)

type ProfileUserHandler struct {
	usecase usecase.ProfileUseCase
}
type ProfileUserRequest struct {
}

type ProfileUserResponse struct {
	User domain.User
}

func NewProfileUserHandler(usecase usecase.ProfileUseCase) *ProfileUserHandler {
	return &ProfileUserHandler{usecase: usecase}
}

func (h *ProfileUserHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *ProfileUserRequest) (*ProfileUserResponse, error) {
	user, err := h.usecase.Execute(fbrCtx, ctx)
	if err != nil {
		return nil, err
	}
	return &ProfileUserResponse{User: *user}, nil
}
