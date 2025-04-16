package fallow

import (
	"context"
	"socialmedia/fallow/app/fallow/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FallowRequestHandler struct {
	usecase usecase.FallowRequestUseCase
}
type FallowRequestRequest struct {
	FollowingID uuid.UUID `json:"following_id"`
}

type FallowRequestResponse struct {
	Message string
}

func NewFallowRequestHandler(usecase usecase.FallowRequestUseCase) *FallowRequestHandler {
	return &FallowRequestHandler{usecase: usecase}
}

func (h *FallowRequestHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *FallowRequestRequest) (*FallowRequestResponse, error) {
	message, err := h.usecase.Execute(fbrCtx, ctx, req.FollowingID)
	if err != nil {
		return nil, err
	}
	return &FallowRequestResponse{Message: message}, nil
}
