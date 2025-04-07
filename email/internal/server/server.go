package server

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Port         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewFiberApp(cfg Config) *fiber.App {
	return fiber.New(fiber.Config{
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Concurrency:  256 * 1024,
	})
}

func Start(app *fiber.App, port string) error {
	return app.Listen(fmt.Sprintf("0.0.0.0:%s", port))
}
