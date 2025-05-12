package chat

import (
	"context"
	"socialmedia/chat/app/chat/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EditMessageContentHandler struct {
	usecase usecase.EditMessageContentUseCase
}

type EditMessageContentRequest struct {
	MessageID uuid.UUID `params:"message_id"`
	Contetnt  string    `json:"content"`
}

type EditMessageContentResponse struct {
	Message string
}

func NewEditMessageContentHandler(usecase usecase.EditMessageContentUseCase) *EditMessageContentHandler {
	return &EditMessageContentHandler{
		usecase: usecase,
	}
}

func (h *EditMessageContentHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *EditMessageContentRequest) (*EditMessageContentResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx, req.MessageID, req.Contetnt)
	if err != nil {
		return nil, err
	}
	return &EditMessageContentResponse{
		Message: "message content edit",
	}, nil
}
