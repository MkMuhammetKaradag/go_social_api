package follow

import (
	"context"
	"socialmedia/follow/app/follow/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AcceptFollowRequestHandler struct {
	usecase usecase.AcceptFollowRequestUseCase
}
type AcceptFollowRequest struct {
	RequestID uuid.UUID `json:"request_id"`
}

type AcceptFollowResponse struct {
	Message string
}

func NewAcceptFollowRequestHandler(usecase usecase.AcceptFollowRequestUseCase) *AcceptFollowRequestHandler {
	return &AcceptFollowRequestHandler{usecase: usecase}
}

func (h *AcceptFollowRequestHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *AcceptFollowRequest) (*AcceptFollowResponse, error) {
	message, err := h.usecase.Execute(fbrCtx, ctx, req.RequestID)
	if err != nil {
		return nil, err
	}
	return &AcceptFollowResponse{Message: message}, nil
}
