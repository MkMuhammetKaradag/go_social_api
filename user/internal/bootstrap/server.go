package bootstrap

import (
	"socialmedia/shared/middlewares"
	"socialmedia/user/internal/handler"
	"socialmedia/user/internal/server"
	"socialmedia/user/pkg/config"
	"time"

	user "socialmedia/user/app/user/handler"

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

	// Korumalı rotalar
	authMiddleware := middlewares.NewAuthMiddleware(redisRepo)
	protected := app.Group("/", authMiddleware.Authenticate())
	{
		profileHandler := httpHandlers["profile"].(*user.ProfileUserHandler)
		updateHandler := httpHandlers["update"].(*user.UpdateUserHandler)
		getUserHandler := httpHandlers["getUser"].(*user.GetUserHandler)

		protected.Get("/profile", handler.HandleWithFiber[user.ProfileUserRequest, user.ProfileUserResponse](profileHandler))
		protected.Post("/update", handler.HandleWithFiber[user.UpdateUserRequest, user.UpdateUserResponse](updateHandler))
		protected.Get("/:id", handler.HandleWithFiber[user.GetUserRequest, user.GetUserResponse](getUserHandler))
	}

	return app
}
