package notification

import (
	"context"

	"socialmedia/notification/app/notification/usecase"

	"github.com/gofiber/fiber/v2"
)

type DeleteAllNotificationsHandler struct {
	usecase usecase.DeleteAllNotificationsUseCase
}

type DeleteAllNotificationsRequest struct {
}

type DeleteAllNotificationsResponse struct {
	Message string
}

func NewDeleteAllNotificationsHandler(usecase usecase.DeleteAllNotificationsUseCase) *DeleteAllNotificationsHandler {
	return &DeleteAllNotificationsHandler{
		usecase: usecase,
	}
}

func (h *DeleteAllNotificationsHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *DeleteAllNotificationsRequest) (*DeleteAllNotificationsResponse, error) {

	err := h.usecase.Execute(fbrCtx, ctx)
	if err != nil {
		return nil, err
	}

	return &DeleteAllNotificationsResponse{
		Message: "All notifications deleted",
	}, nil
}
