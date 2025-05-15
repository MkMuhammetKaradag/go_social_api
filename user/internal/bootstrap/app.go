package bootstrap

import (
	"socialmedia/shared/messaging"
	grpcserver "socialmedia/user/internal/grpc"
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
	userRedisRepo   UserRedisRepository
	myWS            Hub
	rabbitMQ        Messaging
	fiberApp        *fiber.App
	messageHandlers map[messaging.MessageType]MessageHandler
	httpHandlers    map[string]interface{}
	wsHandlers      map[string]interface{}
	grpcSrv         *grpcserver.Server
}

func NewApp(config config.Config) *App {
	grpcSrv := grpcserver.NewGRPCServer()
	app := &App{
		config: config,
	}

	// Bağımlılıkları başlat
	app.initDependencies(grpcSrv)

	return app
}

func (a *App) initDependencies(grpcSrv *grpcserver.Server) {
	// Database ve Redis başlat
	a.grpcSrv = grpcSrv
	a.repo = InitDatabase(a.config)
	a.redisRepo = InitRedis(a.config)

	a.userRedisRepo = InitUserRedis(a.config)
	a.myWS = InitWebsocket(a.userRedisRepo)

	// fmt.Printf("initDependencies: repo address: %p\n", a.repo)

	// Message handler'larını hazırla
	a.messageHandlers = SetupMessageHandlers(a.repo, a.redisRepo)

	// Messaging yapılandırması
	a.rabbitMQ = SetupMessaging(a.messageHandlers, a.config)

	// HTTP handler'larını hazırla
	a.httpHandlers = SetupHTTPHandlers(a.repo, a.redisRepo, a.rabbitMQ)

	a.wsHandlers = SetupWSHandlers(a.repo, a.userRedisRepo, a.myWS)

	// HTTP sunucusu kurulumu
	a.fiberApp = SetupServer(a.config, a.httpHandlers, a.wsHandlers, a.repo, a.redisRepo, a.rabbitMQ)

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
	go a.grpcSrv.Start(":" + a.config.Server.Port)
	zap.L().Info("Server started on port", zap.String("port", a.config.Server.Port))

	// Graceful shutdown için bekle
	graceful.WaitForShutdown(a.fiberApp, 5*time.Second)
}
