package chat

import (
	"context"
	"socialmedia/chat/app/chat/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RenameConversationHandler struct {
	usecase usecase.RenameConversationUseCase
}

type RenameConversationRequest struct {
	ConversationID   uuid.UUID `params:"conversation_id"`
	ConversationName string    `json:"conversation_name"`
}

type RenameConversationResponse struct {
	Message string
}

func NewRenameConversationHandler(usecase usecase.RenameConversationUseCase) *RenameConversationHandler {
	return &RenameConversationHandler{
		usecase: usecase,
	}
}

func (h *RenameConversationHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *RenameConversationRequest) (*RenameConversationResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx, req.ConversationID, req.ConversationName)
	if err != nil {
		return nil, err
	}
	return &RenameConversationResponse{
		Message: "Conversation name changed",
	}, nil
}
