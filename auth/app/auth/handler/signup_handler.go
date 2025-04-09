package auth

import (
	"context"
	"socialmedia/auth/app/auth/usecase"
)

type SignUpAuthRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	// FirstName      string `json:"first_name"`
	// LastName       string `json:"last_name"`
	// Birthdate      string `json:"birthdate"`
	// PhoneNumber    string `json:"phone_number"`
	// ProfilePicture string `json:"profile_picture"`
	// Bio            string `json:"bio"`
}

type SignUpAuthResponse struct {
	Message             string `json:"message"`
	UserActivationToken string `json:"userActivationToken"`
}
type SignUpAuthHandler struct {
	usecase usecase.SignUpUseCase
}

func NewSignUpAuthHandler(usecase usecase.SignUpUseCase) *SignUpAuthHandler {
	return &SignUpAuthHandler{
		usecase: usecase,
	}
}

func (h *SignUpAuthHandler) Handle(ctx context.Context, req *SignUpAuthRequest) (*SignUpAuthResponse, error) {
	activationToken, err := h.usecase.Execute(ctx, &usecase.SignUpRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &SignUpAuthResponse{Message: "User Created", UserActivationToken: *activationToken}, nil
}
