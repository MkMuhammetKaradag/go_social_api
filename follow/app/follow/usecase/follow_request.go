package usecase

import (
	"context"
	"socialmedia/follow/domain"
	"socialmedia/shared/messaging"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type followRequestUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
	rabbitMQ    RabbitMQ
}

func NewFollowRequestUseCase(sessionRepo RedisRepository, repository Repository, rabbitMQ RabbitMQ) FollowRequestUseCase {
	return &followRequestUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
		rabbitMQ:    rabbitMQ,
	}
}

func (u *followRequestUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, followingID uuid.UUID) (string, error) {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return "", domain.ErrNotFoundAuthorization
	}

	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return "", err
	}
	hasBlock, err := u.repository.HasBlockRelationship(ctx, currrentUserID, followingID)
	if err != nil {
		return "", err
	}

	if hasBlock {
		return "", domain.ErrBlockedUser
	}
	isPrivate, err := u.repository.IsPrivate(ctx, followingID)
	if err != nil {
		return "", err
	}

	if isPrivate {
		err = u.repository.CreateFollowRequest(ctx, currrentUserID, followingID)
		if err != nil {
			return "", err
		}
		followMessage := messaging.Message{
			Type:       messaging.UserTypes.FollowRequest,
			ToServices: []messaging.ServiceType{messaging.UserService, messaging.ChatService},
			Data: map[string]interface{}{
				"follower_id":  currrentUserID,
				"following_id": followingID,
				"status":       "pending",
			},
			Critical: false,
		}

		if err := u.rabbitMQ.PublishMessage(ctx, followMessage); err != nil {
			return "", err
		}
		return "Follow request sent", nil
	} else {
		err = u.repository.CreateFollow(ctx, currrentUserID, followingID)
		if err != nil {
			return "", err
		}
		followMessage := messaging.Message{
			Type:       messaging.UserTypes.UserFollowed,
			ToServices: []messaging.ServiceType{messaging.UserService, messaging.ChatService},
			Data: map[string]interface{}{
				"follower_id":  currrentUserID,
				"following_id": followingID,
				"status":       "following",
			},
			Critical: true,
		}

		if err := u.rabbitMQ.PublishMessage(ctx, followMessage); err != nil {
			return "", err
		}
		return "User followed successfully", nil
	}

}
