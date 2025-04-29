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

func SetupServer(config config.Config, httpHandlers map[string]interface{}, wsHandlers map[string]interface{}, repo Repository, redisRepo RedisRepository, rabbitMQ Messaging) *fiber.App {
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
		searchUsersHandler := httpHandlers["searchusers"].(*user.SearchUserHandler)
		updateAvatarHandler := httpHandlers["avatar"].(*user.UpdateAvatarHandler)
		updateBannerHandler := httpHandlers["banner"].(*user.UpdateBannerHandler)

		protected.Get("/profile", handler.HandleWithFiber[user.ProfileUserRequest, user.ProfileUserResponse](profileHandler))
		protected.Post("/update", handler.HandleWithFiber[user.UpdateUserRequest, user.UpdateUserResponse](updateHandler))
		protected.Get("/searchusers", handler.HandleWithFiber[user.SearchUserRequest, user.SearchUserResponse](searchUsersHandler))
		protected.Post("/avatar", handler.HandleWithFiber[user.UpdateAvatarRequest, user.UpdateAvatarResponse](updateAvatarHandler))
		protected.Post("/banner", handler.HandleWithFiber[user.UpdateBannerRequest, user.UpdateBannerResponse](updateBannerHandler))
		protected.Get("/:id", handler.HandleWithFiber[user.GetUserRequest, user.GetUserResponse](getUserHandler))

		wsRoute := app.Group("/ws")
		userStatusPublishHandler := wsHandlers["publishstatus"].(*user.UserStatusPublishHandler)
		wsRoute.Get("/hi", handler.HandleWithFiberWS[user.UserStatusPublishRequest](userStatusPublishHandler))

	}

	return app
}
