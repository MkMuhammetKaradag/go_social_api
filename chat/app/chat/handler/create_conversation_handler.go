package chat

import (
	"context"
	"socialmedia/chat/app/chat/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateConversationHandler struct {
	usecase usecase.CreateConversationUseCase
}
type CreateConversationRequest struct {
	UserIDs []uuid.UUID `json:"user_ids"`
	Name    string      `json:"name,omitempty"`
}

type CreateConversationResponse struct {
	Message string
}

func NewCreateConversationHandler(usecase usecase.CreateConversationUseCase) *CreateConversationHandler {
	return &CreateConversationHandler{usecase: usecase}
}

func (h *CreateConversationHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *CreateConversationRequest) (*CreateConversationResponse, error) {
	isGroup := len(req.UserIDs) > 1 //|| req.Name != ""
	err := h.usecase.Execute(fbrCtx, ctx, req.UserIDs, req.Name, isGroup)
	if err != nil {
		return nil, err
	}
	return &CreateConversationResponse{Message: "chat created"}, nil
}
