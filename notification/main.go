package main

import (
	"socialmedia/notification/internal/bootstrap"
	"socialmedia/notification/pkg/config"

	"go.uber.org/zap"
)

func main() {
	appConfig := config.Read()
	defer zap.L().Sync()
	zap.L().Info("app starting...", zap.String("app name", appConfig.App.Name))

	app := bootstrap.NewApp(appConfig)
	app.Start()
}


