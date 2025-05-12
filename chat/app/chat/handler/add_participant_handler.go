package chat

import (
	"context"
	"socialmedia/chat/app/chat/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AddParticipantHandler struct {
	usecase usecase.AddParticipantUseCase
}

type AddParticipantRequest struct {
	ConversationID uuid.UUID `params:"conversation_id"`
	UserID         uuid.UUID `json:"user_id"`
}

type AddParticipantResponse struct {
	Message string
}

func NewAddParticipantHandler(usecase usecase.AddParticipantUseCase) *AddParticipantHandler {
	return &AddParticipantHandler{
		usecase: usecase,
	}
}

func (h *AddParticipantHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *AddParticipantRequest) (*AddParticipantResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx, req.ConversationID, req.UserID)
	if err != nil {
		return nil, err
	}
	return &AddParticipantResponse{
		Message: "user added chat",
	}, nil
}
