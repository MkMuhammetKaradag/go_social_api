package main

import (
	"log"
	"os"
	"socialmedia/auth/pkg/graceful"
	"socialmedia/email/app/email/handler"
	"socialmedia/email/app/email/usecase"
	"socialmedia/email/infra/mailer"
	"socialmedia/email/internal/consumer"
	"socialmedia/email/internal/server"
	"socialmedia/email/pkg/config"
	"socialmedia/shared/messaging"
	"time"
)

func main() {
	appConfig := config.Read()
	smtpMailer := mailer.NewSMTPMailer(appConfig)
	activationUsecase := usecase.NewActivationService(smtpMailer, "./templates")
	passwordResetUsecase := usecase.NewPasswordResetService(smtpMailer, "./templates")

	activationHandler := handler.NewActivationEmailHandler(activationUsecase)
	passwordResetHandler := handler.NewPasswordResetEmailHandler(passwordResetUsecase)

	messageRouter := func(msg messaging.Message) error {
		switch msg.Type {
		case messaging.EmailTypes.ActivateUser:
			return activationHandler.HandleEmail(msg)
		case messaging.EmailTypes.ForgotPassword:
			return passwordResetHandler.HandleEmail(msg)
		default:
			return nil
		}
	}
	rabbit, err := consumer.StartEmailConsumer(messageRouter)

	if err != nil {
		log.Fatal("RabbitMQ başlatılamadı:", err)
	}
	defer rabbit.Close()

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
