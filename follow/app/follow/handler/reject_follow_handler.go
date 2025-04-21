package follow

import (
	"context"
	"socialmedia/follow/app/follow/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RejectFollowRequestHandler struct {
	usecase usecase.RejectFollowRequestUseCase
}
type RejectFollowRequest struct {
	RequestID uuid.UUID `json:"request_id"`
}

type RejectFollowResponse struct {
	Message string
}

func NewRejectFollowRequestHandler(usecase usecase.RejectFollowRequestUseCase) *RejectFollowRequestHandler {
	return &RejectFollowRequestHandler{usecase: usecase}
}

func (h *RejectFollowRequestHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *RejectFollowRequest) (*RejectFollowResponse, error) {
	message, err := h.usecase.Execute(fbrCtx, ctx, req.RequestID)
	if err != nil {
		return nil, err
	}
	return &RejectFollowResponse{Message: message}, nil
}
