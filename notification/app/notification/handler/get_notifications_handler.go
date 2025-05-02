package notification

import (
	"context"
	"socialmedia/notification/app/notification/usecase"
	"socialmedia/notification/domain"

	"github.com/gofiber/fiber/v2"
)

type GetNotificationsHandler struct {
	usecase usecase.GetNotificationsUseCase
}
type GetNotificationsRequest struct {
	Limit int64 `json:"limit"`
	Skip  int64 `json:"skip,omitempty"`
}

type GetNotificationsResponse struct {
	Notifications []domain.Notification
}

func NewGetNotificationsHandler(usecase usecase.GetNotificationsUseCase) *GetNotificationsHandler {
	return &GetNotificationsHandler{usecase: usecase}
}

func (h *GetNotificationsHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *GetNotificationsRequest) (*GetNotificationsResponse, error) {

	notifications, err := h.usecase.Execute(fbrCtx, ctx, req.Limit, req.Skip)
	if err != nil {
		return nil, err
	}

	return &GetNotificationsResponse{Notifications: notifications}, nil
}
