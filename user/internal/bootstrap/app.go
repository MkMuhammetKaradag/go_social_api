package bootstrap

import (
	"socialmedia/shared/messaging"
	"socialmedia/user/pkg/config"
	"socialmedia/user/pkg/graceful"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type App struct {
	config          config.Config
	repo            Repository
	redisRepo       RedisRepository
	rabbitMQ        Messaging
	fiberApp        *fiber.App
	messageHandlers map[messaging.MessageType]MessageHandler
	httpHandlers    map[string]interface{}
}

func NewApp(config config.Config) *App {
	app := &App{
		config: config,
	}

	// Bağımlılıkları başlat
	app.initDependencies()

	return app
}

func (a *App) initDependencies() {
	// Database ve Redis başlat
	a.repo = InitDatabase(a.config)
	a.redisRepo = InitRedis(a.config)

	// fmt.Printf("initDependencies: repo address: %p\n", a.repo)

	// Message handler'larını hazırla
	a.messageHandlers = SetupMessageHandlers(a.repo, a.redisRepo)

	// Messaging yapılandırması
	a.rabbitMQ = SetupMessaging(a.messageHandlers, a.config)

	// HTTP handler'larını hazırla
	a.httpHandlers = SetupHTTPHandlers(a.repo, a.redisRepo, a.rabbitMQ)

	// HTTP sunucusu kurulumu
	a.fiberApp = SetupServer(a.config, a.httpHandlers, a.repo, a.redisRepo, a.rabbitMQ)
}

func (a *App) Start() {
	// Messaging başlat
	defer a.rabbitMQ.Close()

	// HTTP sunucusunu başlat
	go func() {
		port := a.config.Server.Port
		if err := a.fiberApp.Listen(":" + port); err != nil {
			zap.L().Error("Failed to start server", zap.Error(err))
		}
	}()

	zap.L().Info("Server started on port", zap.String("port", a.config.Server.Port))

	// Graceful shutdown için bekle
	graceful.WaitForShutdown(a.fiberApp, 5*time.Second)
}
