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

	// Korumalı rotalar

	authMiddleware := middlewares.NewAuthMiddleware(redisRepo)
	protected := app.Group("/", authMiddleware.Authenticate())
	{
		protected.Post("/createconversation", handler.HandleWithFiber[chat.CreateConversationRequest, chat.CreateConversationResponse](createConversationHandler))
		protected.Post("/createmessage", handler.HandleWithFiber[chat.CreateMessageRequest, chat.CreateMessageResponse](createMessageHandler))
		protected.Post("/conservation/:conservation_id/add-participant", handler.HandleWithFiber[chat.AddParticipantRequest, chat.AddParticipantResponse](addParticipantHandler))
		protected.Delete("/conservation/:conservation_id/remove-participant", handler.HandleWithFiber[chat.RemoveParticipantRequest, chat.RemoveParticipantResponse](removeParticipantHandler))
		protected.Post("/conservation/:conservation_id/promote-to-admin", handler.HandleWithFiber[chat.PromoteToAdminRequest, chat.PromoteToAdminResponse](promoteToAdminHandler))
		protected.Post("/conservation/:conservation_id/demote-from-admin", handler.HandleWithFiber[chat.DemoteFromAdminRequest, chat.DemoteFromAdminResponse](demoteFromAdminHandler))
		protected.Delete("/messages/:message_id", handler.HandleWithFiber[chat.DeleteMessageRequest, chat.DeleteMessageResponse](deleteMessageHandler))

		wsRoute := app.Group("/ws")
		wsRoute.Get("/message/:chatID", handler.HandleWithFiberWS[chat.ChatWebSocketListenRequest](chatListenHandler))

	}

	return app
}
