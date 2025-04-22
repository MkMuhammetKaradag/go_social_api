package follow

import (
	"context"
	"socialmedia/follow/app/follow/usecase"
	"socialmedia/follow/domain"

	"github.com/gofiber/fiber/v2"
)

type GetBlockedUsersHandler struct {
	usecase usecase.GetBlockedUsersUseCase
}
type GetBlockedUsersRequest struct {
}

type GetBlockedUsersResponse struct {
	Message      string
	BlockedUsers []*domain.BlockedUser `json:"blocked_users"`
}

func NewGetBlockedUsersHandler(usecase usecase.GetBlockedUsersUseCase) *GetBlockedUsersHandler {
	return &GetBlockedUsersHandler{usecase: usecase}
}

func (h *GetBlockedUsersHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *GetBlockedUsersRequest) (*GetBlockedUsersResponse, error) {
	requests, err := h.usecase.Execute(fbrCtx, ctx)
	if err != nil {
		return nil, err
	}
	return &GetBlockedUsersResponse{Message: "User follow request", BlockedUsers: requests}, nil
}
