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

func SetupServer(config config.Config, httpHandlers map[string]interface{}, repo Repository, redisRepo RedisRepository, rabbitMQ Messaging) *fiber.App {
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

	// Korumalı rotalar

	authMiddleware := middlewares.NewAuthMiddleware(redisRepo)
	protected := app.Group("/", authMiddleware.Authenticate())
	{
		// profileHandler := httpHandlers["profile"].(*user.ProfileUserHandler)
		protected.Post("/createconversation", handler.HandleWithFiber[chat.CreateConversationRequest, chat.CreateConversationResponse](createConversationHandler))
		protected.Get("/profile", func(c *fiber.Ctx) error {
			return c.SendString("Hello, World!")
		})

	}

	return app
}
