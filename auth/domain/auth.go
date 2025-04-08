package domain

import "time"

type Auth struct {
	ID               int64     `json:"id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	ActivationCode   string    `json:"activationCode"`
	ActivationExpiry time.Time `json:"activationExpiry"`
}

type AuthResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
