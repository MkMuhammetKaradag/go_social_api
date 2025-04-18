package bootstrap

import (
	auth "socialmedia/auth/app/auth/handler"
	"socialmedia/auth/app/auth/usecase"
	"socialmedia/shared/messaging"
)

func SetupMessageHandlers(repo Repository, redisRepo RedisRepository) map[messaging.MessageType]MessageHandler {
	return map[messaging.MessageType]MessageHandler{}
}

func SetupHTTPHandlers(jwtHelper JwtHelper, repo Repository, redisRepo RedisRepository, rabbitMQ Messaging) map[string]interface{} {
	signUpUseCase := usecase.NewSignUpUseCase(repo, rabbitMQ, jwtHelper)
	forgotPasswordUseCase := usecase.NewForgotPasswordUseCase(repo, rabbitMQ)
	activateUseCase := usecase.NewActivateUseCase(repo, jwtHelper, rabbitMQ)
	signInUseCase := usecase.NewSignInUseCase(repo, redisRepo)
	logoutUseCase := usecase.NewLogoutUseCase(redisRepo)
	resetPasswordUseCase := usecase.NewResetPasswordUseCase(repo, redisRepo)

	signUpAuthHandler := auth.NewSignUpAuthHandler(signUpUseCase)
	forgotPasswordAuthHandler := auth.NewForgotPasswordAuthHandler(forgotPasswordUseCase)
	activateAuthHandler := auth.NewActivateAuthHandler(activateUseCase)
	signInAuthHandler := auth.NewSignInAuthHandler(signInUseCase)
	logoutAuthHandler := auth.NewLogoutAuthHandler(logoutUseCase)
	resetPasswordAuthHandler := auth.NewResetPasswordAuthHandler(resetPasswordUseCase)

	return map[string]interface{}{
		"signup":         signUpAuthHandler,
		"signin":         signInAuthHandler,
		"activate":       activateAuthHandler,
		"forgotpassword": forgotPasswordAuthHandler,
		"resetpassword":  resetPasswordAuthHandler,
		"logout":         logoutAuthHandler,
	}
}
