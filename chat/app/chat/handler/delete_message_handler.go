package chat

import (
	"context"
	"socialmedia/chat/app/chat/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DeleteMessageHandler struct {
	usecase usecase.DeleteMessageUseCase
}

type DeleteMessageRequest struct {
	MessageID uuid.UUID `params:"message_id"`
}

type DeleteMessageResponse struct {
	Message string
}

func NewDeleteMessageHandler(usecase usecase.DeleteMessageUseCase) *DeleteMessageHandler {
	return &DeleteMessageHandler{
		usecase: usecase,
	}
}

func (h *DeleteMessageHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *DeleteMessageRequest) (*DeleteMessageResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx, req.MessageID)
	if err != nil {
		return nil, err
	}
	return &DeleteMessageResponse{
		Message: "message  deleted",
	}, nil
}
