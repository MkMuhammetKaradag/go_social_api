package server

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Config server yapılandırması
type Config struct {
	Port         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewFiberApp yeni bir fiber uygulaması oluşturur
func NewFiberApp(cfg Config) *fiber.App {
	return fiber.New(fiber.Config{
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Concurrency:  256 * 1024,
	})
}

// Start server'ı başlatır
func Start(app *fiber.App, port string) error {
	return app.Listen(fmt.Sprintf("0.0.0.0:%s", port))
}
