package follow

import (
	"context"
	"socialmedia/follow/app/follow/usecase"
	"socialmedia/follow/domain"

	"github.com/gofiber/fiber/v2"
)

type IncomingRequestsHandler struct {
	usecase usecase.IncomingRequestsUseCase
}
type IncomingRequestsRequest struct {
}

type IncomingRequestsResponse struct {
	Message string
	Users   []*domain.User `json:"users"`
}

func NewIncomingRequestsHandler(usecase usecase.IncomingRequestsUseCase) *IncomingRequestsHandler {
	return &IncomingRequestsHandler{usecase: usecase}
}

func (h *IncomingRequestsHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *IncomingRequestsRequest) (*IncomingRequestsResponse, error) {
	users, err := h.usecase.Execute(fbrCtx, ctx)
	if err != nil {
		return nil, err
	}
	return &IncomingRequestsResponse{Message: "User blocked successfully", Users: users}, nil
}
