package chat

import (
	"context"
	"socialmedia/chat/app/chat/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MarkMessagesAsReadHandler struct {
	usecase usecase.MarkMessagesAsReadUseCase
}

type MarkMessagesAsReadRequest struct {
	MessagesIDs []uuid.UUID `json:"messages_ids"`
}

type MarkMessagesAsReadResponse struct {
	Message string
}

func NewMarkMessagesAsReadHandler(usecase usecase.MarkMessagesAsReadUseCase) *MarkMessagesAsReadHandler {
	return &MarkMessagesAsReadHandler{
		usecase: usecase,
	}
}

func (h *MarkMessagesAsReadHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *MarkMessagesAsReadRequest) (*MarkMessagesAsReadResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx, req.MessagesIDs)
	if err != nil {
		return nil, err
	}
	return &MarkMessagesAsReadResponse{
		Message: "user readed messages",
	}, nil
}
