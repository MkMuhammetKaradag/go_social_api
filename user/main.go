package main

import (
	"os"
	"socialmedia/shared/messaging"
	"socialmedia/shared/middlewares"
	follow "socialmedia/user/app/follow/handler"
	followUseCase "socialmedia/user/app/follow/usecase"
	user "socialmedia/user/app/user/handler"
	userUseCase "socialmedia/user/app/user/usecase"
	"socialmedia/user/internal/handler"
	"socialmedia/user/internal/initializer"
	"socialmedia/user/internal/server"
	"socialmedia/user/pkg/config"
	"socialmedia/user/pkg/graceful"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type MessageHandler interface {
	Handle(msg messaging.Message) error
}

func main() {
	appConfig := config.Read()
	defer zap.L().Sync()
	zap.L().Info("app starting...", zap.String("app name", appConfig.App.Name))
	repo := initializer.InitDatabase(appConfig)
	redisRepo := initializer.InitRedis(appConfig)

	followRequestUseCase := followUseCase.NewFollowRequestUseCase(repo)
	followRequestHandler := follow.NewFollowRequestHandler(followRequestUseCase)
	unfollowRequestUseCase := followUseCase.NewUnFollowRequestUseCase(repo)
	unfollowRequestHandler := follow.NewUnFollowRequestHandler(unfollowRequestUseCase)

	blockUserUseCase := followUseCase.NewBlockUserUseCase(repo)
	blockUserHandler := follow.NewBlockUserHandler(blockUserUseCase)
	unblockUserUseCase := followUseCase.NewUnBlockUserUseCase(repo)
	unblockUserHandler := follow.NewUnBlockUserHandler(unblockUserUseCase)

	createUserUseCase := userUseCase.NewCreateUserUseCase(repo)
	createUserHandler := user.NewCreatedUserHandler(createUserUseCase)

	var handlers = map[messaging.MessageType]MessageHandler{
		messaging.UserTypes.UserCreated:     createUserHandler,
		messaging.UserTypes.UserFollowed:    followRequestHandler,
		messaging.UserTypes.FollowRequest:   followRequestHandler,
		messaging.UserTypes.UnFollowRequest: unfollowRequestHandler,
		messaging.UserTypes.UserBlocked:     blockUserHandler,
		messaging.UserTypes.UserUnBlocked:   unblockUserHandler,
	}

	messageRouter := func(msg messaging.Message) error {
		handler, ok := handlers[msg.Type]
		if !ok {
			return nil
		}
		return handler.Handle(msg)
	}

	rabbitMQ := initializer.InitMessaging(messageRouter)
	defer rabbitMQ.Close()

	profileUseCase := userUseCase.NewProfileUseCase(redisRepo, repo)
	updateUserUseCase := userUseCase.NewUpdateUserUseCase(redisRepo, repo)
	getUserUseCase := userUseCase.NewGetUserUseCase(redisRepo, repo)

	profileUserHandler := user.NewProfileUserHandler(profileUseCase)
	updateUserHandler := user.NewUpdateUserHandler(updateUserUseCase)
	getUserHandler := user.NewGetUserHandler(getUserUseCase)

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
		protected.Post("/update", handler.HandleWithFiber[user.UpdateUserRequest, user.UpdateUserResponse](updateUserHandler))
		protected.Get("/:id", handler.HandleWithFiber[user.GetUserRequest, user.GetUserResponse](getUserHandler))

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
