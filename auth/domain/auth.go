package domain

import "time"

type Auth struct {
	ID               string    `json:"id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	ActivationCode   string    `json:"activationCode"`
	ActivationExpiry time.Time `json:"activationExpiry"`
}

type AuthResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ForgotPassword struct {
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}
