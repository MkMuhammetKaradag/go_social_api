package follow

import (
	"context"
	"socialmedia/chat/app/follow/usecase"
	"socialmedia/chat/domain"
	"socialmedia/shared/messaging"

	"github.com/google/uuid"
)

type UnFollowRequestHandler struct {
	usecase usecase.UnFollowRequestUseCase
}

func NewUnFollowRequestHandler(usecase usecase.UnFollowRequestUseCase) *UnFollowRequestHandler {
	return &UnFollowRequestHandler{usecase: usecase}
}

func (h *UnFollowRequestHandler) Handle(msg messaging.Message) error {

	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return domain.ErrInvalidMessageFormat
	}

	UnfollowerIDStr, ok := data["unfollower_id"].(string)
	if !ok {
		return domain.ErrMissingId
	}

	UnfollowingIDStr, ok := data["unfollowing_id"].(string)
	if !ok {
		return domain.ErrMissingId
	}

	UnfollowerID, err := uuid.Parse(UnfollowerIDStr)
	if err != nil {
		return domain.ErrMissingId
	}

	UnfollowingID, err := uuid.Parse(UnfollowingIDStr)
	if err != nil {
		return domain.ErrMissingId
	}

	ctx := context.Background()

	return h.usecase.Execute(ctx, UnfollowerID, UnfollowingID)

}
