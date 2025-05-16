package bootstrap

import (
	auth "socialmedia/auth/app/auth/handler"
	"socialmedia/auth/internal/handler"
	"socialmedia/auth/internal/server"
	"socialmedia/auth/pkg/config"
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

	signUpAuthHandler := httpHandlers["signup"].(*auth.SignUpAuthHandler)
	forgotPasswordAuthHandler := httpHandlers["forgotpassword"].(*auth.ForgotPasswordAuthHandler)
	activateAuthHandler := httpHandlers["activate"].(*auth.ActivateAuthHandler)
	signInAuthHandler := httpHandlers["signin"].(*auth.SignInAuthHandler)
	resetPasswordAuthHandler := httpHandlers["resetpassword"].(*auth.ResetPasswordAuthHandler)

	app.Post("/signup", handler.HandleBasic[auth.SignUpAuthRequest, auth.SignUpAuthResponse](signUpAuthHandler))
	app.Post("/signin", handler.HandleWithFiber[auth.SignInAuthRequest, auth.SignInAuthResponse](signInAuthHandler))
	app.Post("/activate", handler.HandleBasic[auth.ActivateAuthRequest, auth.ActivateAuthResponse](activateAuthHandler))
	app.Post("/forgotpassword", handler.HandleBasic[auth.ForgotPasswordAuthRequest, auth.ForgotPasswordAuthResponse](forgotPasswordAuthHandler))
	app.Post("/resetpassword", handler.HandleBasic[auth.ResetPasswordAuthRequest, auth.ResetPasswordAuthResponse](resetPasswordAuthHandler))

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Korumalı rotalar
	authMiddleware := middlewares.NewAuthMiddleware(redisRepo)
	protected := app.Group("/", authMiddleware.Authenticate())
	{

		logoutAuthHandler := httpHandlers["logout"].(*auth.LogoutAuthHandler)
		allLogoutAuthHandler := httpHandlers["alllogout"].(*auth.AllLogoutAuthHandler)

		protected.Post("/logout", handler.HandleWithFiber[auth.LogoutAuthRequest, auth.LogoutAuthResponse](logoutAuthHandler))
		protected.Post("/all-logout", handler.HandleWithFiber[auth.AllLogoutAuthRequest, auth.AllLogoutAuthResponse](allLogoutAuthHandler))
	}

	return app
}
