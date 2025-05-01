package domain

import "errors"

var (
	ErrInvalidMessageFormat  = errors.New("invalid message format")
	ErrMissingEmail          = errors.New("missing email")
	ErrMissingId             = errors.New("missing id")
	ErrMissingUserName       = errors.New("missing userName")
	ErrNotFoundAuthorization = errors.New("authorization not found ")
	ErrBlockedUser           = errors.New("cannot follow due to block relationship")
)
