package follow

import (
	"context"
	"socialmedia/follow/app/follow/usecase"
	"socialmedia/follow/domain"

	"github.com/gofiber/fiber/v2"
)

type OutgoingRequestsHandler struct {
	usecase usecase.OutgoingRequestsUseCase
}
type OutgoingRequestsRequest struct {
}

type OutgoingRequestsResponse struct {
	Message  string
	Requests []*domain.FollowRequestUser `json:"requests"`
}

func NewOutgoingRequestsHandler(usecase usecase.OutgoingRequestsUseCase) *OutgoingRequestsHandler {
	return &OutgoingRequestsHandler{usecase: usecase}
}

func (h *OutgoingRequestsHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *OutgoingRequestsRequest) (*OutgoingRequestsResponse, error) {
	requests, err := h.usecase.Execute(fbrCtx, ctx)
	if err != nil {
		return nil, err
	}
	return &OutgoingRequestsResponse{Message: "User follow request", Requests: requests}, nil
}
