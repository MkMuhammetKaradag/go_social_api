package main

import (
	"fmt"
	"log"
	"os"
	"socialmedia/auth/pkg/graceful"
	"socialmedia/email/app/email/handler"
	"socialmedia/email/infra/mailer"
	"socialmedia/email/internal/consumer"
	"socialmedia/email/internal/server"
	"socialmedia/email/pkg/config"
	"time"
)

func main() {
	appConfig := config.Read()
	mailer := mailer.NewSMTPMailer(appConfig)
	handler := handler.NewEmailHandler(mailer)
	rabbit, err := consumer.StartEmailConsumer(handler.HandleMessage)
	if err != nil {
		log.Fatal("RabbitMQ başlatılamadı:", err)
	}
	defer rabbit.Close()
	fmt.Println(appConfig.SMTP.Email)

	serverConfig := server.Config{
		Port:         appConfig.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app := server.NewFiberApp(serverConfig)
	go func() {
		if err := server.Start(app, appConfig.Server.Port); err != nil {

			os.Exit(1)
		}
	}()

	graceful.WaitForShutdown(app, 5*time.Second)

}
