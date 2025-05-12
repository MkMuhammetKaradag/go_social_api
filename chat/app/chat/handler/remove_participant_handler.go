package chat

import (
	"context"
	"socialmedia/chat/app/chat/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RemoveParticipantHandler struct {
	usecase usecase.RemoveParticipantUseCase
}

type RemoveParticipantRequest struct {
	ConversationID uuid.UUID `params:"conversation_id"`
	UserID         uuid.UUID `json:"user_id"`
}

type RemoveParticipantResponse struct {
	Message string
}

func NewRemoveParticipantHandler(usecase usecase.RemoveParticipantUseCase) *RemoveParticipantHandler {
	return &RemoveParticipantHandler{
		usecase: usecase,
	}
}

func (h *RemoveParticipantHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *RemoveParticipantRequest) (*RemoveParticipantResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx, req.ConversationID, req.UserID)
	if err != nil {
		return nil, err
	}
	return &RemoveParticipantResponse{
		Message: "user rmoved chat",
	}, nil
}
