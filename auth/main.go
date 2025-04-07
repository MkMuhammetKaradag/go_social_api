// cmd/auth/main.go
package main

import (
	"log"
	"os"
	"socialmedia/auth/app/auth"
	"socialmedia/auth/infra/postgres"
	"socialmedia/auth/infra/redisrepo"
	"socialmedia/auth/internal/handler"
	"socialmedia/auth/internal/server"
	"socialmedia/auth/pkg/config"
	"socialmedia/auth/pkg/graceful"
	"time"

	_ "socialmedia/auth/pkg/log"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	
	appConfig := config.Read()
	defer zap.L().Sync()

	zap.L().Info("app starting...", zap.String("app name", appConfig.App.Name))


	repo, err := postgres.NewPgRepository("postgres://myuser:mypassword@localhost:5432/auth?sslmode=disable")
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	redisRepo, err := redisrepo.NewRedisRepository("localhost:6379", "", 0)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	
	signUpAuthHandler := auth.NewSignUpAuthHandler(repo)
	signInAuthHandler := auth.NewSignInAuthHandler(repo, redisRepo)

	
	serverConfig := server.Config{
		Port:         appConfig.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	
	app := server.NewFiberApp(serverConfig)

	
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/signup", handler.HandleBasic[auth.SignUpAuthRequest, auth.SignUpAuthResponse](signUpAuthHandler))
	app.Post("/signin", handler.HandleWithFiber[auth.SignInAuthRequest, auth.SignInAuthResponse](signInAuthHandler))


	go func() {
		if err := server.Start(app, appConfig.Server.Port); err != nil {
			zap.L().Error("Failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()

	zap.L().Info("Server started on port", zap.String("port", appConfig.Server.Port))

	
	graceful.WaitForShutdown(app, 5*time.Second)
}
