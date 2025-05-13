package chat

import (
	"context"
	"socialmedia/chat/app/chat/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MarkConversationMessagesAsReadHandler struct {
	usecase usecase.MarkConversationMessagesAsReadUseCase
}

type MarkConversationMessagesAsReadRequest struct {
	// conversationID uuid.UUID `params:"conversation_id"`
	ConversationID uuid.UUID `params:"conversation_id"`
}

type MarkConversationMessagesAsReadResponse struct {
	Message string
}

func NewMarkConversationMessagesAsReadHandler(usecase usecase.MarkConversationMessagesAsReadUseCase) *MarkConversationMessagesAsReadHandler {
	return &MarkConversationMessagesAsReadHandler{
		usecase: usecase,
	}
}

func (h *MarkConversationMessagesAsReadHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *MarkConversationMessagesAsReadRequest) (*MarkConversationMessagesAsReadResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx, req.ConversationID)
	if err != nil {
		return nil, err
	}
	return &MarkConversationMessagesAsReadResponse{
		Message: "user readed conversation",
	}, nil
}
