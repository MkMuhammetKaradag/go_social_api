package follow

import (
	"context"
	"socialmedia/shared/messaging"
	"socialmedia/user/app/follow/usecase"
	"socialmedia/user/domain"

	"github.com/google/uuid"
)

type UnBlockUserHandler struct {
	usecase usecase.UnBlockUserUseCase
}

func NewUnBlockUserHandler(usecase usecase.UnBlockUserUseCase) *UnBlockUserHandler {
	return &UnBlockUserHandler{usecase: usecase}
}

func (h *UnBlockUserHandler) Handle(msg messaging.Message) error {

	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return domain.ErrInvalidMessageFormat
	}

	blockerIDStr, ok := data["unblocker_id"].(string)
	if !ok {
		return domain.ErrMissingId
	}

	blockedIDStr, ok := data["unblocked_id"].(string)
	if !ok {
		return domain.ErrMissingId
	}

	blockerID, err := uuid.Parse(blockerIDStr)
	if err != nil {
		return domain.ErrMissingId
	}

	blockedID, err := uuid.Parse(blockedIDStr)
	if err != nil {
		return domain.ErrMissingId
	}

	ctx := context.Background()

	return h.usecase.Execute(ctx, blockerID, blockedID)

}
