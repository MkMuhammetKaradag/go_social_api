package fallow

import (
	"context"
	"socialmedia/fallow/app/fallow/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UnblockUserHandler struct {
	usecase usecase.UnblockUserUseCase
}
type UnblockUserRequest struct {
	BlockedID uuid.UUID `json:"blocked_id"`
}

type UnblockUserResponse struct {
	Message string
} 

func NewUnblockUserHandler(usecase usecase.UnblockUserUseCase) *UnblockUserHandler {
	return &UnblockUserHandler{usecase: usecase}
}

func (h *UnblockUserHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *UnblockUserRequest) (*UnblockUserResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx, req.BlockedID)
	if err != nil {
		return nil, err
	}
	return &UnblockUserResponse{Message: "User unblocked successfully"}, nil
}
