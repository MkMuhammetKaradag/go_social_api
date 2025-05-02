package bootstrap

import (
	notification "socialmedia/notification/app/notification/handler"
	"socialmedia/notification/internal/handler"
	"socialmedia/notification/internal/server"
	"socialmedia/notification/pkg/config"
	"socialmedia/shared/middlewares"

	"time"

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

	// HTTP handler'ları al
	// httpHandlers := SetupHTTPHandlers( repo, redisRepo, rabbitMQ) // Repo parametresi gerekiyorsa düzeltin

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	getNotificationsHandler := httpHandlers["getnotifications"].(*notification.GetNotificationsHandler)
	// Korumalı rotalar
	authMiddleware := middlewares.NewAuthMiddleware(redisRepo)
	protected := app.Group("/", authMiddleware.Authenticate())
	{

		protected.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("Hello, World-2!")
		})
		protected.Get("/notifications", handler.HandleWithFiber[notification.GetNotificationsRequest, notification.GetNotificationsResponse](getNotificationsHandler))

	}

	return app
}
