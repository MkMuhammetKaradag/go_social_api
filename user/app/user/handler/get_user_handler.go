package user

import (
	"context"
	"socialmedia/user/app/user/usecase"
	"socialmedia/user/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetUserHandler struct {
	usecase usecase.GetUserUseCase
}
type GetUserRequest struct {
}

type GetUserResponse struct {
	User *domain.User
}

func NewGetUserHandler(usecase usecase.GetUserUseCase) *GetUserHandler {
	return &GetUserHandler{usecase: usecase}
}

func (h *GetUserHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	id := fbrCtx.Params("id")

	targetUserID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	user, err := h.usecase.Execute(fbrCtx, ctx, targetUserID)
	if err != nil {
		return nil, err
	}
	return &GetUserResponse{User: user}, nil
}
