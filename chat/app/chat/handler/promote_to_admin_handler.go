package chat

import (
	"context"
	"socialmedia/chat/app/chat/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PromoteToAdminHandler struct {
	usecase usecase.PromoteToAdminUseCase
}

type PromoteToAdminRequest struct {
	ConversationID uuid.UUID `params:"conservation_id"`
	UserID         uuid.UUID `json:"user_id"`
}

type PromoteToAdminResponse struct {
	Message string
}

func NewPromoteToAdminHandler(usecase usecase.PromoteToAdminUseCase) *PromoteToAdminHandler {
	return &PromoteToAdminHandler{
		usecase: usecase,
	}
}

func (h *PromoteToAdminHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *PromoteToAdminRequest) (*PromoteToAdminResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx, req.ConversationID, req.UserID)
	if err != nil {
		return nil, err
	}
	return &PromoteToAdminResponse{
		Message: "User became admin ",
	}, nil
}
