package follow

import (
	"context"
	"socialmedia/follow/app/follow/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BlockUserHandler struct {
	usecase usecase.BlockUserUseCase
}
type BlockUserRequest struct {
	BlockedID uuid.UUID `json:"blocked_id"`
}

type BlockUserResponse struct {
	Message string
}

func NewBlockUserHandler(usecase usecase.BlockUserUseCase) *BlockUserHandler {
	return &BlockUserHandler{usecase: usecase}
}

func (h *BlockUserHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *BlockUserRequest) (*BlockUserResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx, req.BlockedID)
	if err != nil {
		return nil, err
	}
	return &BlockUserResponse{Message: "User blocked successfully"}, nil
}
