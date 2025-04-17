package follow

import (
	"context"
	"socialmedia/shared/messaging"
	"socialmedia/user/app/follow/usecase"
	"socialmedia/user/domain"

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

	UnfollowerIDStr, ok := data["Unfollower_id"].(string)
	if !ok {
		return domain.ErrMissingEmail
	}

	UnfollowingIDStr, ok := data["Unfollowing_id"].(string)
	if !ok {
		return domain.ErrMissingId
	}

	status, ok := data["status"].(string)
	if !ok {
		return domain.ErrMissingUserName
	}

	UnfollowerID, err := uuid.Parse(UnfollowerIDStr)
	if err != nil {
		return domain.ErrMissingEmail
	}

	UnfollowingID, err := uuid.Parse(UnfollowingIDStr)
	if err != nil {
		return domain.ErrMissingId
	}

	ctx := context.Background()

	return h.usecase.Execute(ctx, UnfollowerID, UnfollowingID, status)

}
