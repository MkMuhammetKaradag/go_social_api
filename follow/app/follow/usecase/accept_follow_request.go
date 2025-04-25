package usecase

import (
	"context"
	"socialmedia/follow/domain"
	"socialmedia/shared/messaging"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type acceptFollowRequestUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
	rabbitMQ    RabbitMQ
}

func NewAcceptFollowRequestUseCase(sessionRepo RedisRepository, repository Repository, rabbitMQ RabbitMQ) AcceptFollowRequestUseCase {
	return &acceptFollowRequestUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
		rabbitMQ:    rabbitMQ,
	}
}

func (u *acceptFollowRequestUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, requestID uuid.UUID) (string, error) {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return "", domain.ErrNotFoundAuthorization
	}

	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return "", err
	}

	followerID, err := u.repository.AcceptFollowRequest(ctx, requestID, currrentUserID)
	if err != nil {
		return "", err
	}

	followMessage := messaging.Message{
		Type:       messaging.UserTypes.UserFollowed,
		ToServices: []messaging.ServiceType{messaging.UserService, messaging.ChatService},
		Data: map[string]interface{}{
			"follower_id":  followerID,
			"following_id": currrentUserID,
			"status":       "following",
		},
		Critical: true,
	}
	if err := u.rabbitMQ.PublishMessage(ctx, followMessage); err != nil {
		return "", err
	}
	return "Follow request  accepted", nil

}
