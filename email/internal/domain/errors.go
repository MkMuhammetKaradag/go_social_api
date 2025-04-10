package domain

import "errors"

var (
	ErrInvalidMessageFormat  = errors.New("invalid message format")
	ErrMissingEmail          = errors.New("missing email")
	ErrMissingActivationCode = errors.New("missing activation code")
	ErrMissingResetLink      = errors.New("missing reset link")
	ErrMissingTemplateName   = errors.New("missing template name")
)
