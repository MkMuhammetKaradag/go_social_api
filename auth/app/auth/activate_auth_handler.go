package auth

import (
	"context"
	"errors"
	"fmt"
	"socialmedia/auth/domain"

	"github.com/gofiber/fiber/v2"
)

type ActivateAuthRequest struct {
	ActivationToken string `json:"activationToken" binding:"required"`
	ActivationCode  string `json:"activationCode" binding:"required,min=4,max=4"`
}

type ActivateAuthResponse struct {
	Auth *domain.AuthResponse `json:"auth"`
}

type ActivateAuthHandler struct {
	repository Repository
	jwtHelper  JwtHelper
}

func NewActivateAuthHandler(repository Repository, jwtHelper JwtHelper) *ActivateAuthHandler {
	return &ActivateAuthHandler{
		repository: repository,
		jwtHelper:  jwtHelper,
	}
}

func (h *ActivateAuthHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *ActivateAuthRequest) (*ActivateAuthResponse, error) {

	claims, err := h.jwtHelper.VerifyToken(req.ActivationToken)
	if err != nil {
		return nil, fmt.Errorf("error verifying token: %w", err)
	}
	activationCode := claims["activationCode"].(string)
	if activationCode != req.ActivationCode {
		return nil, errors.New("activation code mismatch")
	}
	userEmail, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid user data in token")
	}

	auth, err := h.repository.Activate(ctx, userEmail, activationCode)
	if err != nil {
		return nil, err
	}

	response := &domain.AuthResponse{
		ID:       auth.ID,
		Username: auth.Username,
		Email:    auth.Email,
	}
	return &ActivateAuthResponse{Auth: response}, nil
}
