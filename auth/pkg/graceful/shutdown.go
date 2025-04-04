package graceful

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// WaitForShutdown uygulama kapatma sinyallerini bekler
func WaitForShutdown(app *fiber.App, timeout time.Duration) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	zap.L().Info("Shutting down server...")

	if err := app.ShutdownWithTimeout(timeout); err != nil {
		zap.L().Error("Error during server shutdown", zap.Error(err))
	}

	zap.L().Info("Server gracefully stopped")
}
