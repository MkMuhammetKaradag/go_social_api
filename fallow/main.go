package main

import (
	"os"
	user "socialmedia/fallow/app/user/handler"
	userUseCase "socialmedia/fallow/app/user/usecase"
	"socialmedia/fallow/internal/initializer"
	"socialmedia/fallow/internal/server"
	"socialmedia/fallow/pkg/config"
	"socialmedia/fallow/pkg/graceful"
	"socialmedia/shared/messaging"
	"socialmedia/shared/middlewares"

	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.Read()
	defer zap.L().Sync()
	zap.L().Info("app starting...", zap.String("app name", appConfig.App.Name))
	repo := initializer.InitDatabase(appConfig)
	redisRepo := initializer.InitRedis(appConfig)

	createUserUseCase := userUseCase.NewCreateUserUseCase(repo)
	createUserHandler := user.NewCreatedUserHandler(createUserUseCase)
	messageRouter := func(msg messaging.Message) error {
		switch msg.Type {
		case messaging.UserTypes.UserCreated:
			err := createUserHandler.Handle(msg)
			return err

			return nil

		default:
			return nil
		}
	}

	rabbitMQ := initializer.InitMessaging(messageRouter)
	defer rabbitMQ.Close()

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
		protected.Get("/pro", func(c *fiber.Ctx) error {
			return c.SendString("Hello, World!-protecdet")
		})

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
