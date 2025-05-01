package usecase

import (
	"context"
	"fmt"
	"socialmedia/chat/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type createConversationUseCase struct {
	repository Repository
}

func NewCreateConversationUseCase(repository Repository) CreateConversationUseCase {
	return &createConversationUseCase{
		repository: repository,
	}
}

func (u *createConversationUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, userIDs []uuid.UUID, name string, isGroup bool) error {

	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}

	currentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}

	// UUID'leri benzersiz hale getir
	uniqueUserIDMap := make(map[uuid.UUID]struct{})
	for _, id := range userIDs {
		uniqueUserIDMap[id] = struct{}{}
	}
	uniqueUserIDMap[currentUserID] = struct{}{} // current user'ı da ekle

	// Benzersiz ID listesini oluştur
	uniqueUserIDs := make([]uuid.UUID, 0, len(uniqueUserIDMap))
	for id := range uniqueUserIDMap {
		isblock, _ := u.repository.HasBlockRelationship(ctx, currentUserID, id)
		if isblock {
			continue
		}
		uniqueUserIDs = append(uniqueUserIDs, id)
	}

	// 2 veya daha az katılımcı varsa grup değildir
	isGroup = len(uniqueUserIDs) > 2

	conversation, blockedParticipant, err := u.repository.CreateConversation(ctx, currentUserID, isGroup, name, uniqueUserIDs)
	if err != nil {
		return err

	}
	if blockedParticipant != nil && len(*blockedParticipant) > 0 {
		fmt.Println("conversation ID:", conversation.ID, "block partispant:", blockedParticipant)
		fmt.Println("bildirim yollandı")
	}

	return nil
}
