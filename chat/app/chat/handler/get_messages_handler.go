package chat

import (
	"context"
	"socialmedia/chat/app/chat/usecase"
	"socialmedia/chat/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetMessagesHandler struct {
	usecase usecase.GetMessagesUseCase
}

type GetMessagesRequest struct {
	ConversationID uuid.UUID `params:"conversation_id"`
	Limit     int64     `json:"limit"`
	Skip      int64     `json:"skip,omitempty"`
}

type GetMessagesResponse struct {
	Messages []domain.Message
}

func NewGetMessagesHandler(usecase usecase.GetMessagesUseCase) *GetMessagesHandler {
	return &GetMessagesHandler{
		usecase: usecase,
	}
}

func (h *GetMessagesHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *GetMessagesRequest) (*GetMessagesResponse, error) {
	messages,err := h.usecase.Execute(fbrCtx, ctx, req.ConversationID, req.Limit,req.Skip)
	if err != nil {
		return nil, err
	}
	return &GetMessagesResponse{
		Messages: messages,
	}, nil
}
