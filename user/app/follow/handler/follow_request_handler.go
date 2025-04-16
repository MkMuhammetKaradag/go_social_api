package follow

import (
	"context"
	"socialmedia/shared/messaging"
	"socialmedia/user/app/follow/usecase"
	"socialmedia/user/domain"

	"github.com/google/uuid"
)

type FollowRequestHandler struct {
	usecase usecase.FollowRequestUseCase
}

func NewFollowRequestHandler(usecase usecase.FollowRequestUseCase) *FollowRequestHandler {
	return &FollowRequestHandler{usecase: usecase}
}

func (h *FollowRequestHandler) Handle(msg messaging.Message) error {

	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return domain.ErrInvalidMessageFormat
	}

	followerIDStr, ok := data["follower_id"].(string)
	if !ok {
		return domain.ErrMissingEmail
	}

	followingIDStr, ok := data["following_id"].(string)
	if !ok {
		return domain.ErrMissingId
	}

	status, ok := data["status"].(string)
	if !ok {
		return domain.ErrMissingUserName
	}

	followerID, err := uuid.Parse(followerIDStr)
	if err != nil {
		return domain.ErrMissingEmail
	}

	followingID, err := uuid.Parse(followingIDStr)
	if err != nil {
		return domain.ErrMissingId
	}

	ctx := context.Background()

	return h.usecase.Execute(ctx, followerID, followingID, status)

}
