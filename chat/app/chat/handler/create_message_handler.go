package chat

import (
	"context"
	"socialmedia/chat/app/chat/usecase"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateMessageHandler struct {
	usecase usecase.CreateMessageUseCase
}

type CreateMessageRequest struct {
	ConversationID  uuid.UUID `json:"conversation_id"`
	Content         string    `json:"content"`
	AttachmentURLs  []string  `json:"attachment_urls,omitempty"`
	AttachmentTypes []string  `json:"attachment_types,omitempty"`
}

type CreateMessageResponse struct {
	MessageID uuid.UUID `json:"message_id"`
	Message   string    `json:"message"`
}

func NewCreateMessageHandler(usecase usecase.CreateMessageUseCase) *CreateMessageHandler {
	return &CreateMessageHandler{usecase: usecase}
}

func (h *CreateMessageHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *CreateMessageRequest) (*CreateMessageResponse, error) {
	userID, err := getUserIDFromContext(fbrCtx)
	if err != nil {
		return nil, err
	}

	messageID, err := h.usecase.Execute(ctx, req.ConversationID, userID, req.Content, req.AttachmentURLs, req.AttachmentTypes)
	if err != nil {
		return nil, err
	}

	return &CreateMessageResponse{
		MessageID: messageID,
		Message:   "message sent successfully",
	}, nil
}

func getUserIDFromContext(ctx *fiber.Ctx) (uuid.UUID, error) {

	userData, _ := middlewares.GetUserData(ctx)
	return uuid.Parse(userData["id"])
}
