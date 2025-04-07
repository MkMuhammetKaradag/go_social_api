package auth

import (
	"context"
	"fmt"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
)

type LogoutAuthRequest struct {
}

type LogoutAuthResponse struct {
	Message string `json:"message"`
}

type LogoutAuthHandler struct {
	sessionRepo RedisRepository
}

func NewLogoutAuthHandler(sessionRepo RedisRepository) *LogoutAuthHandler {
	return &LogoutAuthHandler{
		sessionRepo: sessionRepo,
	}
}

func (h *LogoutAuthHandler) Handle(fbrCtx *fiber.Ctx, ctx context.Context, req *LogoutAuthRequest) (*LogoutAuthResponse, error) {
	cookieSessionId := fbrCtx.Cookies("session_id")
	if cookieSessionId == "" {
		return nil, fmt.Errorf("user not logout")
	}

	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return nil, fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}

	userID := userData["id"]
	if userID == "" {
		return nil, fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}

	if err := h.sessionRepo.DeleteSession(ctx, cookieSessionId, userID); err != nil {
		return nil, err
	}

	fbrCtx.ClearCookie("session_id")

	return &LogoutAuthResponse{Message: "logout user "}, nil
}
