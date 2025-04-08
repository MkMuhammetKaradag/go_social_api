package jwthelper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtHelperService struct {
	secretKey string
}

func NewJwtHelperService(secretKey string) *JwtHelperService {
	return &JwtHelperService{secretKey: secretKey}
}

func (j *JwtHelperService) SignToken(payload jwt.MapClaims, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(expiration).Unix(),
	}
	for key, value := range payload {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JwtHelperService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
