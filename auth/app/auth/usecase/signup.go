package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"socialmedia/auth/domain"

	"time"

	"github.com/golang-jwt/jwt"
)

type signUpUseCase struct {
	repository Repository
	rabbitMQ   RabbitMQ
	jwtHelper  JwtHelper
}

type SignUpRequest struct {
	Username string
	Email    string
	Password string
}

func NewSignUpUseCase(repository Repository, rabbitMQ RabbitMQ, jwtHelper JwtHelper) SignUpUseCase {
	return &signUpUseCase{
		repository: repository,
		rabbitMQ:   rabbitMQ,
		jwtHelper:  jwtHelper,
	}
}

func (u *signUpUseCase) Execute(ctx context.Context, req *SignUpRequest) (*string, error) {

	activationCode := generateActivationCode()

	auth := &domain.Auth{
		Username:         req.Username,
		Email:            req.Email,
		Password:         req.Password,
		ActivationCode:   activationCode,
		ActivationExpiry: time.Now().Add(5 * time.Minute),
	}

	err := u.repository.SignUp(ctx, auth)
	if err != nil {
		return nil, err
	}

	payload := jwt.MapClaims{
		"email": req.Email,
	}

	activationToken, err := u.jwtHelper.SignToken(payload, 10*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("jwt sign error: %w", err)
	}
	fmt.Println(activationCode)
	// emailMessage := messaging.Message{
	// 	Type:       messaging.EmailTypes.ActivateUser,
	// 	ToServices: []messaging.ServiceType{messaging.EmailService},
	// 	Data: map[string]interface{}{
	// 		"email":           req.Email,
	// 		"activation_code": activationCode,
	// 		"template_name":   "activation_email.html",
	// 		"userName":        req.Username,
	// 	},
	// 	Critical: false,
	// }

	// if err := u.rabbitMQ.PublishMessage(ctx, emailMessage); err != nil {
	// 	return nil, err
	// }

	return &activationToken, nil
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateActivationCode() string {
	return fmt.Sprintf("%04d", r.Intn(10000))
}
