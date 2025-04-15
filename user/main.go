package main

import (
	"os"
	"socialmedia/shared/messaging"
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
	redisRepo := initializer.InitRedis(appConfig)

	createUserUseCase := usecase.NewCreateUserUseCase(repo)
	createUserHandler := user.NewCreatedUserHandler(createUserUseCase)

	messageRouter := func(msg messaging.Message) error {

		// fmt.Println("mesage geldi main:", msg.Type)
		switch msg.Type {
		case messaging.UserTypes.UserCreated:
			// fmt.Println("case gitrdi created")
			err := createUserHandler.Handle(msg)

			return err

		default:
			return nil
		}
	}

	rabbitMQ := initializer.InitMessaging(messageRouter)
	defer rabbitMQ.Close()

	profileUseCase := usecase.NewProfileUseCase(redisRepo, repo)
	updateUserUseCase := usecase.NewUpdateUserUseCase(redisRepo, repo)
	getUserUseCase := usecase.NewGetUserUseCase(redisRepo, repo)

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
