package main

import (
	"os"
	follow "socialmedia/follow/app/follow/handler"
	followUseCase "socialmedia/follow/app/follow/usecase"
	user "socialmedia/follow/app/user/handler"
	userUseCase "socialmedia/follow/app/user/usecase"
	"socialmedia/follow/internal/handler"
	"socialmedia/follow/internal/initializer"
	"socialmedia/follow/internal/server"
	"socialmedia/follow/pkg/config"
	"socialmedia/follow/pkg/graceful"
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

	followRequestUseCase := followUseCase.NewFollowRequestUseCase(redisRepo, repo, rabbitMQ)
	blockUserUseCase := followUseCase.NewBlockUserUseCase(redisRepo, repo)
	unblockUserUseCase := followUseCase.NewUnblockUserUseCase(redisRepo, repo)
	fallawRequestHandler := follow.NewFollowRequestHandler(followRequestUseCase)
	blockUserHandler := follow.NewBlockUserHandler(blockUserUseCase)
	unblockUserHandler := follow.NewUnblockUserHandler(unblockUserUseCase)

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
		protected.Post("/follow", handler.HandleWithFiber[follow.FollowRequestRequest, follow.FollowRequestResponse](fallawRequestHandler))
		protected.Post("/block", handler.HandleWithFiber[follow.BlockUserRequest, follow.BlockUserResponse](blockUserHandler))
		protected.Post("/unblock", handler.HandleWithFiber[follow.UnblockUserRequest, follow.UnblockUserResponse](unblockUserHandler))

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
