package follow

import (
	"context"
	"socialmedia/follow/app/follow/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UnFollowRequestHandler struct {
	usecase usecase.UnFollowRequestUseCase
}
type UnFollowRequestRequest struct {
	UnFollowingID uuid.UUID `json:"unfollowing_id"`
}

type UnFollowRequestResponse struct {
	Message string
}

func NewUnFollowRequestHandler(usecase usecase.UnFollowRequestUseCase) *UnFollowRequestHandler {
	return &UnFollowRequestHandler{usecase: usecase}
}

func (h *UnFollowRequestHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *UnFollowRequestRequest) (*UnFollowRequestResponse, error) {
	message, err := h.usecase.Execute(fbrCtx, ctx, req.UnFollowingID)
	if err != nil {
		return nil, err
	}
	return &UnFollowRequestResponse{Message: message}, nil
}
