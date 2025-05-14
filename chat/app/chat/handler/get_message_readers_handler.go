package chat

import (
	"context"
	"socialmedia/chat/app/chat/usecase"
	"socialmedia/chat/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetMessageReadersHandler struct {
	usecase usecase.GetMessageReadersUseCase
}

type GetMessageReadersRequest struct {
	MessageID uuid.UUID `params:"message_id"`
}

type GetMessageReadersResponse struct {
	Users []domain.User
}

func NewGetMessageReadersHandler(usecase usecase.GetMessageReadersUseCase) *GetMessageReadersHandler {
	return &GetMessageReadersHandler{
		usecase: usecase,
	}
}

func (h *GetMessageReadersHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *GetMessageReadersRequest) (*GetMessageReadersResponse, error) {
	users, err := h.usecase.Execute(fbrCtx, ctx, req.MessageID)
	if err != nil {
		return nil, err
	}
	return &GetMessageReadersResponse{
		Users: users,
	}, nil
}
