package domain

import "errors"

var (
	ErrInvalidMessageFormat = errors.New("invalid message format")
	ErrMissingEmail         = errors.New("missing email")
	ErrMissingId            = errors.New("missing id")
	ErrMissingUserName      = errors.New("missing userName")
)
