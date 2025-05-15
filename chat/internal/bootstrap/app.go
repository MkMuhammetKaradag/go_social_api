package bootstrap

import (
	"context"
	"socialmedia/chat/pkg/config"
	"socialmedia/chat/pkg/graceful"
	"socialmedia/shared/messaging"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type App struct {
	config          config.Config
	repo            Repository
	userClient      UserClient
	myWS            Hub
	redisRepo       RedisRepository
	chatRedisRepo   ChatRedisRepository
	rabbitMQ        Messaging
	fiberApp        *fiber.App
	messageHandlers map[messaging.MessageType]MessageHandler
	httpHandlers    map[string]interface{}
	wsHandlers      map[string]interface{}
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
	a.userClient = InitUserClient(a.config)
	a.chatRedisRepo = InitChatRedis(a.config)
	redisClient := a.chatRedisRepo.GetRedisClient()
	ctx := context.Background()
	a.myWS = InitWebsocket(ctx, redisClient, a.repo)

	// fmt.Printf("initDependencies: repo address: %p\n", a.repo)

	// Message handler'larını hazırla
	a.messageHandlers = SetupMessageHandlers(a.repo, a.redisRepo)

	// Messaging yapılandırması
	a.rabbitMQ = SetupMessaging(a.messageHandlers, a.config)

	// HTTP handler'larını hazırla
	a.httpHandlers = SetupHTTPHandlers(a.repo, a.redisRepo, a.chatRedisRepo, a.rabbitMQ,a.userClient)
	a.wsHandlers = SetupWSHandlers(a.repo, a.chatRedisRepo, a.myWS)

	// HTTP sunucusu kurulumu
	a.fiberApp = SetupServer(a.config, a.httpHandlers, a.wsHandlers, a.redisRepo)
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
