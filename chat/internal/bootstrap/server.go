package bootstrap

import (
	chat "socialmedia/chat/app/chat/handler"
	"socialmedia/chat/internal/handler"
	"socialmedia/chat/internal/server"
	"socialmedia/chat/pkg/config"
	"socialmedia/shared/middlewares"
	"time"

	// user "socialmedia/user/app/user/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupServer(config config.Config, httpHandlers map[string]interface{}, wsHandlers map[string]interface{}, redisRepo RedisRepository) *fiber.App {
	serverConfig := server.Config{
		Port:         config.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app := server.NewFiberApp(serverConfig)

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	createConversationHandler := httpHandlers["createconversation"].(*chat.CreateConversationHandler)
	createMessageHandler := httpHandlers["createmessage"].(*chat.CreateMessageHandler)
	chatListenHandler := wsHandlers["chatlisten"].(*chat.ChatWebSocketListenHandler)
	addParticipantHandler := httpHandlers["addparticipant"].(*chat.AddParticipantHandler)
	removeParticipantHandler := httpHandlers["removeparticipant"].(*chat.RemoveParticipantHandler)
	promoteToAdminHandler := httpHandlers["promotetoadmin"].(*chat.PromoteToAdminHandler)
	demoteFromAdminHandler := httpHandlers["demotefromadmin"].(*chat.DemoteFromAdminHandler)
	deleteMessageHandler := httpHandlers["deletemessage"].(*chat.DeleteMessageHandler)
	renameConversationHandler := httpHandlers["renameconversation"].(*chat.RenameConversationHandler)
	editMessageContentHandler := httpHandlers["editmessagecontent"].(*chat.EditMessageContentHandler)
	markMessagesAsReadHandler := httpHandlers["markmessagesasread"].(*chat.MarkMessagesAsReadHandler)
	markConversationMessagesAsReadHandler := httpHandlers["markconversationmessagesasread"].(*chat.MarkConversationMessagesAsReadHandler)
	getMessagesHandler := httpHandlers["getmessages"].(*chat.GetMessagesHandler)
	getMessageReadersHandler := httpHandlers["getmessagereaders"].(*chat.GetMessageReadersHandler)
	deleteAllMessagesFromConversationHandler := httpHandlers["deleteallmessagesfromconversation"].(*chat.DeleteAllMessagesFromConversationHandler)

	// Korumalı rotalar

	authMiddleware := middlewares.NewAuthMiddleware(redisRepo)
	protectedConversation := app.Group("/conversations", authMiddleware.Authenticate())
	{
		protectedConversation.Post("/", handler.HandleWithFiber[chat.CreateConversationRequest, chat.CreateConversationResponse](createConversationHandler))
		protectedConversation.Patch("/:conversation_id/rename", handler.HandleWithFiber[chat.RenameConversationRequest, chat.RenameConversationResponse](renameConversationHandler))
		protectedConversation.Post("/:conversation_id/add-participant", handler.HandleWithFiber[chat.AddParticipantRequest, chat.AddParticipantResponse](addParticipantHandler))
		protectedConversation.Delete("/:conversation_id/remove-participant", handler.HandleWithFiber[chat.RemoveParticipantRequest, chat.RemoveParticipantResponse](removeParticipantHandler))
		protectedConversation.Delete("/:conversation_id/messages", handler.HandleWithFiber[chat.DeleteAllMessagesFromConversationRequest, chat.DeleteAllMessagesFromConversationResponse](deleteAllMessagesFromConversationHandler))
		protectedConversation.Post("/:conversation_id/promote-to-admin", handler.HandleWithFiber[chat.PromoteToAdminRequest, chat.PromoteToAdminResponse](promoteToAdminHandler))
		protectedConversation.Post("/:conversation_id/demote-from-admin", handler.HandleWithFiber[chat.DemoteFromAdminRequest, chat.DemoteFromAdminResponse](demoteFromAdminHandler))
		protectedConversation.Post("/:conversation_id/read", handler.HandleWithFiber[chat.MarkConversationMessagesAsReadRequest, chat.MarkConversationMessagesAsReadResponse](markConversationMessagesAsReadHandler))

	}
	protectedMessage := app.Group("/messages", authMiddleware.Authenticate())
	{
		protectedMessage.Post("/", handler.HandleWithFiber[chat.CreateMessageRequest, chat.CreateMessageResponse](createMessageHandler))
		protectedMessage.Get("/:conversation_id", handler.HandleWithFiber[chat.GetMessagesRequest, chat.GetMessagesResponse](getMessagesHandler))
		protectedMessage.Get("/:message_id/readers", handler.HandleWithFiber[chat.GetMessageReadersRequest, chat.GetMessageReadersResponse](getMessageReadersHandler))
		protectedMessage.Delete("/:message_id", handler.HandleWithFiber[chat.DeleteMessageRequest, chat.DeleteMessageResponse](deleteMessageHandler))
		protectedMessage.Patch("/:message_id/edit-message-content", handler.HandleWithFiber[chat.EditMessageContentRequest, chat.EditMessageContentResponse](editMessageContentHandler))
		protectedMessage.Post("/read", handler.HandleWithFiber[chat.MarkMessagesAsReadRequest, chat.MarkMessagesAsReadResponse](markMessagesAsReadHandler))

	}
	wsRoute := app.Group("/ws")
	wsRoute.Get("/message/:chatID", handler.HandleWithFiberWS[chat.ChatWebSocketListenRequest](chatListenHandler))

	return app
}
