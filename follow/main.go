package main

import (
	"socialmedia/follow/internal/bootstrap"
	"socialmedia/follow/pkg/config"

	"go.uber.org/zap"
)

func main() {
	appConfig := config.Read()
	defer zap.L().Sync()
	zap.L().Info("app starting...", zap.String("app name", appConfig.App.Name))

	app := bootstrap.NewApp(appConfig)
	app.Start()
}
