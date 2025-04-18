package main

import (
	"fmt"
	"socialmedia/auth/internal/bootstrap"
	"socialmedia/auth/pkg/config"

	"go.uber.org/zap"
)

func main() {
	appConfig := config.Read()
	defer zap.L().Sync()
	zap.L().Info("app starting...", zap.String("app name", appConfig.App.Name))

	app := bootstrap.NewApp(appConfig)
	fmt.Println("mami star")

	app.Start()
}
