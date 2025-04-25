package user

import (
	"context"
	"socialmedia/chat/app/user/usecase"
	"socialmedia/chat/domain"
	"socialmedia/shared/messaging"

	"github.com/google/uuid"
)

type UpdateUserHandler struct {
	usecase usecase.UpdateUserUseCase
}

func NewUpdatedUserHandler(UpdatedUserUsecase usecase.UpdateUserUseCase) *UpdateUserHandler {
	return &UpdateUserHandler{
		usecase: UpdatedUserUsecase,
	}
}

func (h *UpdateUserHandler) Handle(msg messaging.Message) error {

	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return domain.ErrInvalidMessageFormat
	}

	userIDStr, ok := data["user_id"].(string)
	if !ok {
		return domain.ErrMissingId
	}
	var userName *string
	if val, ok := data["avatar_url"]; ok {
		if strVal, ok := val.(string); ok && strVal != "" {
			userName = &strVal
		}
	}

	var avatarURL *string
	if val, ok := data["avatar_url"]; ok {
		if strVal, ok := val.(string); ok && strVal != "" {
			avatarURL = &strVal
		}
	}

	var isPrivate *bool
	if val, ok := data["is_private"]; ok {
		if boolVal, ok := val.(bool); ok {
			isPrivate = &boolVal
		}
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return domain.ErrMissingId
	}
	ctx := context.Background()
	return h.usecase.Execute(ctx, userID, userName, avatarURL, isPrivate)
}
