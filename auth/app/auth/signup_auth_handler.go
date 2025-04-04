package auth

import (
	"context"
	"socialmedia/auth/domain"
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
}

func NewSignUpAuthHandler(repository Repository) *SignUpAuthHandler {
	return &SignUpAuthHandler{
		repository: repository,
	}
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

	return &SignUpAuthResponse{Message: "User Created"}, nil
}
