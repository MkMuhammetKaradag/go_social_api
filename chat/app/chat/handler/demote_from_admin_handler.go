package chat

import (
	"context"
	"socialmedia/chat/app/chat/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DemoteFromAdminHandler struct {
	usecase usecase.DemoteFromAdminUseCase
}

type DemoteFromAdminRequest struct {
	ConversationID uuid.UUID `params:"conservation_id"`
	UserID         uuid.UUID `json:"user_id"`
}

type DemoteFromAdminResponse struct {
	Message string
}

func NewDemoteFromAdminHandler(usecase usecase.DemoteFromAdminUseCase) *DemoteFromAdminHandler {
	return &DemoteFromAdminHandler{
		usecase: usecase,
	}
}

func (h *DemoteFromAdminHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *DemoteFromAdminRequest) (*DemoteFromAdminResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx, req.ConversationID, req.UserID)
	if err != nil {
		return nil, err
	}
	return &DemoteFromAdminResponse{
		Message: "user admin was demoted ",
	}, nil
}
