package user

import (
	"context"
	"socialmedia/shared/messaging"
	"socialmedia/user/app/user/usecase"
	"socialmedia/user/domain"
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

	email, ok := data["email"].(string)
	if !ok {
		return domain.ErrMissingEmail
	}

	id, ok := data["id"].(string)
	if !ok {
		return domain.ErrMissingId
	}

	userName, ok := data["username"].(string)
	if !ok {
		return domain.ErrMissingUserName
	}
	ctx := context.Background()
	return h.usecase.Execute(ctx, id, userName, email)
}
