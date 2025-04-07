package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func respondWithError(c *fiber.Ctx, code int, message string) error {
	c.Set("Content-Type", "application/json")
	return c.Status(code).JSON(fiber.Map{"error": message})
}

type AuthMiddleware struct {
	sessionRepo RedisRepository
}

func NewAuthMiddleware(redisRepo RedisRepository) *AuthMiddleware {
	return &AuthMiddleware{sessionRepo: redisRepo}
}

func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {

		var token string

		if strings.Contains(c.Get("Connection"), "Upgrade") && c.Get("Upgrade") == "websocket" {
			token = c.Query("token")
			if token == "" {
				token = c.Get("session_id")
			}
		} else {
			cookieSessionId := c.Cookies("session_id")
			if cookieSessionId == "" {
				return respondWithError(c, fiber.StatusUnauthorized, "Unauthorized: missing session")
			}
			token = cookieSessionId
		}

		userData, err := m.sessionRepo.GetSession(c.UserContext(), token)
		if err != nil {
			return respondWithError(c, fiber.StatusUnauthorized, "missing session")
		}

		c.Locals("userData", userData)

		return c.Next()
	}
}

func GetUserData(c *fiber.Ctx) (map[string]string, bool) {
	userData, ok := c.Locals("userData").(map[string]string)
	return userData, ok
}
