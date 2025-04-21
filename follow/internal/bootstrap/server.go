package bootstrap

import (
	follow "socialmedia/follow/app/follow/handler"
	"socialmedia/follow/internal/handler"
	"socialmedia/follow/internal/server"
	"socialmedia/follow/pkg/config"
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

	// Korumalı rotalar
	authMiddleware := middlewares.NewAuthMiddleware(redisRepo)
	protected := app.Group("/", authMiddleware.Authenticate())
	{
		follawRequestHandler := httpHandlers["follow"].(*follow.FollowRequestHandler)
		unfollawRequestHandler := httpHandlers["unfollow"].(*follow.UnFollowRequestHandler)
		blockUserHandler := httpHandlers["block"].(*follow.BlockUserHandler)
		unblockUserHandler := httpHandlers["unblock"].(*follow.UnblockUserHandler)
		incomingRequestHandler := httpHandlers["incomingRequest"].(*follow.IncomingRequestsHandler)

		protected.Post("/follow", handler.HandleWithFiber[follow.FollowRequestRequest, follow.FollowRequestResponse](follawRequestHandler))
		protected.Post("/unfollow", handler.HandleWithFiber[follow.UnFollowRequestRequest, follow.UnFollowRequestResponse](unfollawRequestHandler))
		
		protected.Post("/block", handler.HandleWithFiber[follow.BlockUserRequest, follow.BlockUserResponse](blockUserHandler))
		protected.Post("/unblock", handler.HandleWithFiber[follow.UnblockUserRequest, follow.UnblockUserResponse](unblockUserHandler))

		protected.Get("/follow/requests/incoming", handler.HandleWithFiber[follow.IncomingRequestsRequest, follow.IncomingRequestsResponse](incomingRequestHandler))

	}

	return app
}
