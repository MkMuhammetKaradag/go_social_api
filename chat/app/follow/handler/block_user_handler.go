package follow

import (
	"context"
	"socialmedia/shared/messaging"
	"socialmedia/chat/app/follow/usecase"
	"socialmedia/chat/domain"

	"github.com/google/uuid"
)

type BlockUserHandler struct {
	usecase usecase.BlockUserUseCase
}

func NewBlockUserHandler(usecase usecase.BlockUserUseCase) *BlockUserHandler {
	return &BlockUserHandler{usecase: usecase}
}

func (h *BlockUserHandler) Handle(msg messaging.Message) error {

	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return domain.ErrInvalidMessageFormat
	}

	blockerIDStr, ok := data["blocker_id"].(string)
	if !ok {
		return domain.ErrMissingId
	}

	blockedIDStr, ok := data["blocked_id"].(string)
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
