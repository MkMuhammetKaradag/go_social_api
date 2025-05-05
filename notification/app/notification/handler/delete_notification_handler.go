package notification

import (
	"context"
	"socialmedia/notification/app/notification/usecase"

	"github.com/gofiber/fiber/v2"
)

type DeleteNotificationHandeler struct {
	usecase usecase.DeleteNotificationUseCase
}

type DeleteNotificationRequest struct {
	NotificationID string `params:"notification_id"`
}
type DeleteNotificationResponse struct {
	Message string
}

func NewDeleteNotificationHandler(usecase usecase.DeleteNotificationUseCase) *DeleteNotificationHandeler {
	return &DeleteNotificationHandeler{
		usecase: usecase,
	}
}

func (h *DeleteNotificationHandeler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *DeleteNotificationRequest) (*DeleteNotificationResponse, error) {
	err := h.usecase.Execute(fbrCtx, ctx, req.NotificationID)
	if err != nil {
		return nil, err
	}
	return &DeleteNotificationResponse{
		Message: "notification deleted",
	}, nil
}
