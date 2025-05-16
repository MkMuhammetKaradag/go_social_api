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



	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	getNotificationsHandler := httpHandlers["getnotifications"].(*notification.GetNotificationsHandler)
	getUnreadNotificationsHandler := httpHandlers["getunreadnotifications"].(*notification.GetUnreadNotificationsHandler)
	markNotificationHandler := httpHandlers["marknotification"].(*notification.MarkNotificationHandler)
	deleteNotificationHandler := httpHandlers["deletenotification"].(*notification.DeleteNotificationHandeler)
	readAllNotificationsHandler := httpHandlers["readallnotifications"].(*notification.ReadAllNotificationsHandler)
	deleteAllNotificationsHandler := httpHandlers["deleteallnotifications"].(*notification.DeleteAllNotificationsHandler)

	// Korumalı rotalar
	authMiddleware := middlewares.NewAuthMiddleware(redisRepo)
	protected := app.Group("/", authMiddleware.Authenticate())
	{

		protected.Get("/notification", handler.HandleWithFiber[notification.GetNotificationsRequest, notification.GetNotificationsResponse](getNotificationsHandler))
		protected.Get("/notification/unread", handler.HandleWithFiber[notification.GetUnreadNotificationsRequest, notification.GetUnreadNotificationsResponse](getUnreadNotificationsHandler))
		protected.Patch("/notification/:notification_id/read", handler.HandleWithFiber[notification.MarkNotificationRequest, notification.MarkNotificationResponse](markNotificationHandler))

		protected.Patch("/notification/read-all", handler.HandleWithFiber[notification.ReadAllNotificationsRequest, notification.ReadAllNotificationsResponse](readAllNotificationsHandler))

		protected.Delete("/notification", handler.HandleWithFiber[notification.DeleteAllNotificationsRequest, notification.DeleteAllNotificationsResponse](deleteAllNotificationsHandler))

		protected.Delete("/notification/:notification_id", handler.HandleWithFiber[notification.DeleteNotificationRequest, notification.DeleteNotificationResponse](deleteNotificationHandler))

	}

	return app
}
