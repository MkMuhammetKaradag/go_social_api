package auth

import (
	"context"
	"fmt"
	"math/rand"
	"socialmedia/auth/domain"
	"socialmedia/shared/messaging"
	"time"
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
	Message string `json:"message"`
}

type SignUpAuthHandler struct {
	repository Repository
	rabbitMQ   RabbitMQ
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewSignUpAuthHandler(repository Repository, rabbitMQ RabbitMQ) *SignUpAuthHandler {
	return &SignUpAuthHandler{
		repository: repository,
		rabbitMQ:   rabbitMQ,
	}
}
func GenerateActivationCode() string {

	num := r.Intn(10000)

	return fmt.Sprintf("%04d", num)

}

func (h *SignUpAuthHandler) Handle(ctx context.Context, req *SignUpAuthRequest) (*SignUpAuthResponse, error) {
	// authId := uuid.New().String()

	auth := &domain.Auth{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.repository.SignUp(ctx, auth)
	if err != nil {
		return nil, err
	}
	activationCode := GenerateActivationCode()
	payload := map[string]interface{}{
		"activationCode": activationCode,
		"user":           req,
	}
	fmt.Println(payload)

	emailMessage := messaging.Message{
		Type:      "active_user",
		ToService: messaging.EmailService,
		Data: map[string]interface{}{
			"email":           req.Email,
			"activation_code": activationCode,
			"template_name":   "activation_email.html",
			"userName":        req.Username,
		},
	}

	if err := h.rabbitMQ.PublishMessage(context.Background(), emailMessage); err != nil {

		return nil, fmt.Errorf("Aktivasyon e-postası gönderilemedi")
	}

	return &SignUpAuthResponse{Message: "User Created"}, nil
}
