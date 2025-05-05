package notification

import (
	"context"

	"socialmedia/notification/app/notification/usecase"

	"github.com/gofiber/fiber/v2"
)

type ReadAllNotificationsHandler struct {
	usecase usecase.ReadAllNotificationsUseCase
}

type ReadAllNotificationsRequest struct {
}
type ReadAllNotificationsResponse struct {
	Message string
}

func NewReadAllNotificationsHandler(usecase usecase.ReadAllNotificationsUseCase) *ReadAllNotificationsHandler {

	return &ReadAllNotificationsHandler{
		usecase: usecase,
	}
}

func (h *ReadAllNotificationsHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *ReadAllNotificationsRequest) (*ReadAllNotificationsResponse, error) {

	err := h.usecase.Execute(fbrCtx, ctx)
	if err != nil {
		return nil, err
	}

	return &ReadAllNotificationsResponse{
		Message: "notification read all",
	}, nil
}
