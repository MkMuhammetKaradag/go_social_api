package usecase

import (
	"context"
	"fmt"
	"socialmedia/follow/domain"
	"socialmedia/shared/messaging"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type unfollowRequestUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
	rabbitMQ    RabbitMQ
}

func NewUnFollowRequestUseCase(sessionRepo RedisRepository, repository Repository, rabbitMQ RabbitMQ) UnFollowRequestUseCase {
	return &unfollowRequestUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
		rabbitMQ:    rabbitMQ,
	}
}

func (u *unfollowRequestUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, followingID uuid.UUID) (string, error) {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return "", domain.ErrNotFoundAuthorization
	}

	currentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return "", fmt.Errorf("invalid user ID: %w", err)
	}

	isFollowing, err := u.repository.IsFollowing(ctx, currentUserID, followingID)
	if err != nil {
		return "", fmt.Errorf("failed to check follow status: %w", err)
	}

	var messageType messaging.MessageType
	var messageText string

	if isFollowing {

		if err := u.repository.DeleteFollow(ctx, currentUserID, followingID); err != nil {
			return "", fmt.Errorf("failed to delete follow: %w", err)
		}
		messageType = messaging.UserTypes.UnFollowRequest
		messageText = "User unfollowed successfully"
	} else {

		if err := u.repository.DeleteFollowRequest(ctx, currentUserID, followingID); err != nil {
			return "", fmt.Errorf("failed to delete follow request: %w", err)
		}
		messageType = messaging.UserTypes.UnFollowRequest
		messageText = "Follow request deleted"
	}

	unfollowMessage := messaging.Message{
		Type:       messageType,
		ToServices: []messaging.ServiceType{messaging.UserService},
		Data: map[string]interface{}{
			"unfollower_id":  currentUserID,
			"unfollowing_id": followingID,
		},
		Critical: true,
	}

	if err := u.rabbitMQ.PublishMessage(ctx, unfollowMessage); err != nil {
		return "", fmt.Errorf("failed to publish unfollow message: %w", err)
	}
	return messageText, nil
}
