// cmd/auth/main.go
package main

import (
	"os"

	auth "socialmedia/auth/app/auth/handler"
	"socialmedia/auth/app/auth/usecase"
	"socialmedia/auth/internal/handler"
	"socialmedia/auth/internal/initializer"
	"socialmedia/auth/internal/server"
	"socialmedia/auth/pkg/config"
	"socialmedia/auth/pkg/graceful"
	"socialmedia/shared/middlewares"
	"time"

	_ "socialmedia/auth/pkg/log"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.Read()
	defer zap.L().Sync()
	zap.L().Info("app starting...", zap.String("app name", appConfig.App.Name))

	repo := initializer.InitDatabase(appConfig)
	redisRepo := initializer.InitRedis(appConfig)
	rabbitMQ := initializer.InitMessaging()
	jwtHelper := initializer.InitJwtHelper(appConfig)

	defer rabbitMQ.Close()

	signUpUseCase := usecase.NewSignUpUseCase(repo, rabbitMQ, jwtHelper)
	forgotPasswordUseCase := usecase.NewForgotPasswordUseCase(repo, rabbitMQ)
	activateUseCase := usecase.NewActivateUseCase(repo, jwtHelper)
	signInUseCase := usecase.NewSignInUseCase(repo, redisRepo)
	logoutUseCase := usecase.NewLogoutUseCase(redisRepo)

	signUpAuthHandler := auth.NewSignUpAuthHandler(signUpUseCase)
	forgotPasswordAuthHandler := auth.NewForgotPasswordAuthHandler(forgotPasswordUseCase)
	activateAuthHandler := auth.NewActivateAuthHandler(activateUseCase)
	signInAuthHandler := auth.NewSignInAuthHandler(signInUseCase)
	logoutAuthHandler := auth.NewLogoutAuthHandler(logoutUseCase)

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
	app.Post("/signup", handler.HandleBasic[auth.SignUpAuthRequest, auth.SignUpAuthResponse](signUpAuthHandler))
	app.Post("/signin", handler.HandleWithFiber[auth.SignInAuthRequest, auth.SignInAuthResponse](signInAuthHandler))
	app.Post("/activate", handler.HandleBasic[auth.ActivateAuthRequest, auth.ActivateAuthResponse](activateAuthHandler))
	app.Post("/forgotpassword", handler.HandleBasic[auth.ForgotPasswordAuthRequest, auth.ForgotPasswordAuthResponse](forgotPasswordAuthHandler))

	protected := app.Group("/", authMiddleware.Authenticate())
	{
		protected.Get("/profile", profileHandler)
		protected.Post("/logout", handler.HandleWithFiber[auth.LogoutAuthRequest, auth.LogoutAuthResponse](logoutAuthHandler))

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

func profileHandler(c *fiber.Ctx) error {
	userData, ok := middlewares.GetUserData(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Kullanıcı bilgisi bulunamadı")
	}
	return c.JSON(userData)
}
