package bootstrap

import (
	chat "socialmedia/chat/app/chat/handler"
	chatUseCase "socialmedia/chat/app/chat/usecase"
	follow "socialmedia/chat/app/follow/handler"
	followUseCase "socialmedia/chat/app/follow/usecase"
	user "socialmedia/chat/app/user/handler"
	userUseCase "socialmedia/chat/app/user/usecase"
	"socialmedia/shared/messaging"
)

func SetupMessageHandlers(repo Repository, redisRepo RedisRepository) map[messaging.MessageType]MessageHandler {
	// Follow related use cases and handlers
	followRequestUseCase := followUseCase.NewFollowRequestUseCase(repo)
	followRequestHandler := follow.NewFollowRequestHandler(followRequestUseCase)
	unfollowRequestUseCase := followUseCase.NewUnFollowRequestUseCase(repo)
	unfollowRequestHandler := follow.NewUnFollowRequestHandler(unfollowRequestUseCase)

	blockUserUseCase := followUseCase.NewBlockUserUseCase(repo)
	blockUserHandler := follow.NewBlockUserHandler(blockUserUseCase)
	unblockUserUseCase := followUseCase.NewUnBlockUserUseCase(repo)
	unblockUserHandler := follow.NewUnBlockUserHandler(unblockUserUseCase)

	// User related use cases and handlers
	createUserUseCase := userUseCase.NewCreateUserUseCase(repo)
	updateUserUseCase := userUseCase.NewUpdateUserUseCase(repo)
	createUserHandler := user.NewCreatedUserHandler(createUserUseCase)
	updateUserHandler := user.NewUpdatedUserHandler(updateUserUseCase)

	return map[messaging.MessageType]MessageHandler{
		messaging.UserTypes.UserCreated:     createUserHandler,
		messaging.UserTypes.UserUpdated:     updateUserHandler,
		messaging.UserTypes.UserFollowed:    followRequestHandler,
		messaging.UserTypes.FollowRequest:   followRequestHandler,
		messaging.UserTypes.UnFollowRequest: unfollowRequestHandler,
		messaging.UserTypes.UserBlocked:     blockUserHandler,
		messaging.UserTypes.UserUnBlocked:   unblockUserHandler,
	}
}

func SetupHTTPHandlers(repo Repository, redisRepo RedisRepository, chatRedisRepo ChatRedisRepository, rabbitMQ Messaging, userClient UserClient) map[string]interface{} {
	createConversationUseCase := chatUseCase.NewCreateConversationUseCase(repo, rabbitMQ)
	createMessageUseCase := chatUseCase.NewCreateMessageUseCase(repo, chatRedisRepo)
	addParticipantUseCase := chatUseCase.NewAddParticipantUseCase(repo, chatRedisRepo)
	removeParticipantUseCase := chatUseCase.NewRemoveParticipantUseCase(repo, chatRedisRepo)
	promoteToAdminUseCase := chatUseCase.NewPromoteToAdminUseCase(repo)
	demoteFromAdminUseCase := chatUseCase.NewDemoteFromAdminUseCase(repo)
	deleteMessageUseCase := chatUseCase.NewDeleteMessageUseCase(repo, chatRedisRepo)
	renameConversationUseCase := chatUseCase.NewRenameConversationUseCase(repo)
	editMessageContentUseCase := chatUseCase.NewEditMessageContentUseCase(repo, chatRedisRepo)
	markMessagesAsReadUseCase := chatUseCase.NewMarkMessagesAsReadUseCase(repo)
	markConversationMessagesAsReadUseCase := chatUseCase.NewMarkConversationMessagesAsReadUseCase(repo)
	getMessagesUseCase := chatUseCase.NewGetMessagesUseCase(repo)
	getMessageReaderssUseCase := chatUseCase.NewGetMessageReadersUseCase(repo)
	deleteAllMessagesFromConversationUseCase := chatUseCase.NewDeleteAllMessagesFromConversationUseCase(repo, chatRedisRepo,userClient)

	createConversationHandler := chat.NewCreateConversationHandler(createConversationUseCase)
	createMessageHandler := chat.NewCreateMessageHandler(createMessageUseCase)
	addParticipantHandler := chat.NewAddParticipantHandler(addParticipantUseCase)
	removeParticipantHandler := chat.NewRemoveParticipantHandler(removeParticipantUseCase)
	promoteToAdminHandler := chat.NewPromoteToAdminHandler(promoteToAdminUseCase)
	demoteFromAdminHandler := chat.NewDemoteFromAdminHandler(demoteFromAdminUseCase)
	deleteMessageHandler := chat.NewDeleteMessageHandler(deleteMessageUseCase)
	renameConversationHandler := chat.NewRenameConversationHandler(renameConversationUseCase)
	editMessageContentHandler := chat.NewEditMessageContentHandler(editMessageContentUseCase)
	markMessagesAsReadHandler := chat.NewMarkMessagesAsReadHandler(markMessagesAsReadUseCase)
	markConversationMessagesAsReadHandler := chat.NewMarkConversationMessagesAsReadHandler(markConversationMessagesAsReadUseCase)
	getMessagesHandler := chat.NewGetMessagesHandler(getMessagesUseCase)
	getMessageReadersHandler := chat.NewGetMessageReadersHandler(getMessageReaderssUseCase)
	deleteAllMessagesFromConversationHandler := chat.NewDeleteAllMessagesFromConversationHandler(deleteAllMessagesFromConversationUseCase)

	return map[string]interface{}{

		"createconversation":                createConversationHandler,
		"createmessage":                     createMessageHandler,
		"addparticipant":                    addParticipantHandler,
		"removeparticipant":                 removeParticipantHandler,
		"promotetoadmin":                    promoteToAdminHandler,
		"demotefromadmin":                   demoteFromAdminHandler,
		"deletemessage":                     deleteMessageHandler,
		"renameconversation":                renameConversationHandler,
		"editmessagecontent":                editMessageContentHandler,
		"markmessagesasread":                markMessagesAsReadHandler,
		"markconversationmessagesasread":    markConversationMessagesAsReadHandler,
		"getmessages":                       getMessagesHandler,
		"getmessagereaders":                 getMessageReadersHandler,
		"deleteallmessagesfromconversation": deleteAllMessagesFromConversationHandler,
	}
}
func SetupWSHandlers(repo Repository, chatRedisRepo ChatRedisRepository, wsHub Hub) map[string]interface{} {
	chatListenUseCase := chatUseCase.NewChatWebSocketListenUseCase(repo, wsHub)

	chatListenHandler := chat.NewChatWebSocketListenHandler(chatListenUseCase)

	return map[string]interface{}{

		"chatlisten": chatListenHandler,
	}
}
