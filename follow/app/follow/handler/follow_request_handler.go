package follow

import (
	"context"
	"socialmedia/follow/app/follow/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FollowRequestHandler struct {
	usecase usecase.FollowRequestUseCase
}
type FollowRequestRequest struct {
	FollowingID uuid.UUID `json:"following_id"`
}

type FollowRequestResponse struct {
	Message string
}

func NewFollowRequestHandler(usecase usecase.FollowRequestUseCase) *FollowRequestHandler {
	return &FollowRequestHandler{usecase: usecase}
}

func (h *FollowRequestHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *FollowRequestRequest) (*FollowRequestResponse, error) {
	message, err := h.usecase.Execute(fbrCtx, ctx, req.FollowingID)
	if err != nil {
		return nil, err
	}
	return &FollowRequestResponse{Message: message}, nil
}
