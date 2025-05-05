package notification

import (
	"context"
	"socialmedia/notification/app/notification/usecase"
	"socialmedia/notification/domain"

	"github.com/gofiber/fiber/v2"
)

type GetUnreadNotificationsHandler struct {
	usecase usecase.GetUnreadNotificationsUseCase
}

type GetUnreadNotificationsRequest struct {
	Limit int64 `json:"limit"`
	Skip  int64 `json:"skip,omitemty"`
}

type GetUnreadNotificationsResponse struct {
	Notifications []domain.Notification
}

func NewGetUnreadNotificationsHandler(usecase usecase.GetUnreadNotificationsUseCase) *GetUnreadNotificationsHandler {
	return &GetUnreadNotificationsHandler{
		usecase: usecase,
	}
}

func (h *GetUnreadNotificationsHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *GetUnreadNotificationsRequest) (*GetUnreadNotificationsResponse, error) {
	notifications, err := h.usecase.Execute(fbrCtx, ctx, req.Limit, req.Skip)

	if err != nil {
		return nil, err
	}
	return &GetUnreadNotificationsResponse{
			Notifications: notifications,
		},
		nil
}
