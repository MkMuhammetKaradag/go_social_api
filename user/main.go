package main

import (
	"fmt"
	"os"
	"socialmedia/shared/middlewares"
	user "socialmedia/user/app/user/handler"
	"socialmedia/user/app/user/usecase"
	"socialmedia/user/internal/handler"
	"socialmedia/user/internal/initializer"
	"socialmedia/user/internal/server"
	"socialmedia/user/pkg/config"
	"socialmedia/user/pkg/graceful"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.Read()
	defer zap.L().Sync()
	zap.L().Info("app starting...", zap.String("app name", appConfig.App.Name))
	repo := initializer.InitDatabase(appConfig)

	fmt.Println(repo)
	redisRepo := initializer.InitRedis(appConfig)
	rabbitMQ := initializer.InitMessaging()
	defer rabbitMQ.Close()

	profileUseCase := usecase.NewProfileUseCase(redisRepo)
	profileUserHandler := user.NewProfileUserHandler(profileUseCase)

	serverConfig := server.Config{
		Port:         appConfig.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app := server.NewFiberApp(serverConfig)
	authMiddleware := middlewares.NewAuthMiddleware(redisRepo)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	protected := app.Group("/", authMiddleware.Authenticate())
	{
		protected.Get("/profile", handler.HandleWithFiber[user.ProfileUserRequest, user.ProfileUserResponse](profileUserHandler))

	}

	go func() {
		if err := server.Start(app, appConfig.Server.Port); err != nil {
			zap.L().Error("Failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()
	zap.L().Info("Server started on port", zap.String("port", appConfig.Server.Port))
	graceful.WaitForShutdown(app, 5*time.Second)
}
