package bootstrap

import (
	"socialmedia/notification/pkg/config"
	"socialmedia/notification/pkg/graceful"
	"socialmedia/shared/messaging"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type App struct {
	config          config.Config
	repo            Repository
	repoMongo       RepositoryMongo
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
	a.repoMongo = InitDatabaseMongo(a.config)
	a.redisRepo = InitRedis(a.config)

	// fmt.Printf("initDependencies: repo address: %p\n", a.repo)

	// Message handler'larını hazırla
	a.messageHandlers = SetupMessageHandlers(a.repo, a.repoMongo, a.redisRepo)

	// Messaging yapılandırması
	a.rabbitMQ = SetupMessaging(a.messageHandlers, a.config)

	// HTTP handler'larını hazırla
	a.httpHandlers = SetupHTTPHandlers(a.repo, a.repoMongo, a.redisRepo, a.rabbitMQ)

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
