package user

import (
	"context"
	"socialmedia/notification/app/user/usecase"
	"socialmedia/notification/domain"
	"socialmedia/shared/messaging"

	"github.com/google/uuid"
)

type CreatedUserHandler struct {
	usecase usecase.CreateUserUseCase
}

func NewCreatedUserHandler(createdUserUsecase usecase.CreateUserUseCase) *CreatedUserHandler {
	return &CreatedUserHandler{
		usecase: createdUserUsecase,
	}
}

func (h *CreatedUserHandler) Handle(msg messaging.Message) error {

	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return domain.ErrInvalidMessageFormat
	}

	idStr, ok := data["id"].(string)
	if !ok {
		return domain.ErrMissingId
	}
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}

	userName, ok := data["username"].(string)
	if !ok {
		return domain.ErrMissingUserName
	}
	ctx := context.Background()
	return h.usecase.Execute(ctx, userID, userName)
}
