package notification

import (
	"context"
	"socialmedia/notification/app/notification/usecase"

	"github.com/gofiber/fiber/v2"
)

type MarkNotificationHandler struct {
	usecase usecase.MarkNotificationUseCase
}
type MarkNotificationRequest struct {
	NotificationID string `params:"notification_id"`
}

type MarkNotificationResponse struct {
	Message string
}

func NewMarkNotificationHandler(usecase usecase.MarkNotificationUseCase) *MarkNotificationHandler {
	return &MarkNotificationHandler{usecase: usecase}
}

func (h *MarkNotificationHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *MarkNotificationRequest) (*MarkNotificationResponse, error) {

	err := h.usecase.Execute(fbrCtx, ctx, req.NotificationID)
	if err != nil {
		return nil, err
	}

	return &MarkNotificationResponse{Message: "notification marked"}, nil
}
