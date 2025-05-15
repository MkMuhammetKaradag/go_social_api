package chat

import (
	"context"
	"socialmedia/chat/app/chat/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DeleteAllMessagesFromConversationHandler struct {
	usecase usecase.DeleteAllMessagesFromConversationUseCase
}

type DeleteAllMessagesFromConversationRequest struct {
	ConversationID uuid.UUID `params:"conversation_id"`
}

type DeleteAllMessagesFromConversationResponse struct {
	Message string
}

func NewDeleteAllMessagesFromConversationHandler(usecase usecase.DeleteAllMessagesFromConversationUseCase) *DeleteAllMessagesFromConversationHandler {
	return &DeleteAllMessagesFromConversationHandler{
		usecase: usecase,
	}
}

func (h *DeleteAllMessagesFromConversationHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *DeleteAllMessagesFromConversationRequest) (*DeleteAllMessagesFromConversationResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx, req.ConversationID)
	if err != nil {
		return nil, err
	}
	return &DeleteAllMessagesFromConversationResponse{
		Message: "all messages  deleted",
	}, nil
}
