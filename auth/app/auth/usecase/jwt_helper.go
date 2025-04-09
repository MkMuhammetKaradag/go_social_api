package usecase

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtHelper interface {
	SignToken(payload jwt.MapClaims, expiration time.Duration) (string, error)
	VerifyToken(tokenStr string) (jwt.MapClaims, error)
}
