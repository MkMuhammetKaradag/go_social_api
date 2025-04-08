package auth

import (
	"context"
	"fmt"
	"math/rand"
	"socialmedia/auth/domain"
	"socialmedia/shared/messaging"
	"time"

	"github.com/golang-jwt/jwt"
)

type SignUpAuthRequest struct {
	Username       string `json:"username" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=8"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Birthdate      string `json:"birthdate"`
	PhoneNumber    string `json:"phone_number"`
	ProfilePicture string `json:"profile_picture"`
	Bio            string `json:"bio"`
}

type SignUpAuthResponse struct {
	Message             string `json:"message"`
	UserActivationToken string `json:"userActivationToken"`
}

type SignUpAuthHandler struct {
	repository Repository
	rabbitMQ   RabbitMQ
	jwtHelper  JwtHelper
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewSignUpAuthHandler(repository Repository, rabbitMQ RabbitMQ, jwtHelper JwtHelper) *SignUpAuthHandler {
	return &SignUpAuthHandler{
		repository: repository,
		rabbitMQ:   rabbitMQ,
		jwtHelper:  jwtHelper,
	}
}
func GenerateActivationCode() string {

	num := r.Intn(10000)

	return fmt.Sprintf("%04d", num)

}

func (h *SignUpAuthHandler) Handle(ctx context.Context, req *SignUpAuthRequest) (*SignUpAuthResponse, error) {
	activationCode := GenerateActivationCode()
	auth := &domain.Auth{
		Username:         req.Username,
		Email:            req.Email,
		Password:         req.Password,
		ActivationCode:   activationCode,
		ActivationExpiry: time.Now().Add(5 * time.Minute),
	}

	err := h.repository.SignUp(ctx, auth)
	if err != nil {
		return nil, err
	}
	fmt.Println(activationCode)
	payload := jwt.MapClaims{
		"activationCode": activationCode,
		"email":          req.Email,
	}
	activationToken, err := h.jwtHelper.SignToken(payload, 10*time.Minute)
	if err != nil {
		fmt.Println("err jwt  token:", err)
	}

	emailMessage := messaging.Message{
		Type:      messaging.EmailTypes.ActivateUser,
		ToService: messaging.EmailService,
		Data: map[string]interface{}{
			"email":           req.Email,
			"activation_code": activationCode,
			"template_name":   "activation_email.html",
			"userName":        req.Username,
		},
	}

	if err := h.rabbitMQ.PublishMessage(context.Background(), emailMessage); err != nil {

		return nil, err
	}

	return &SignUpAuthResponse{Message: "User Created", UserActivationToken: activationToken}, nil
}
