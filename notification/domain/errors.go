package domain

import "errors"

var (
	ErrInvalidMessageFormat  = errors.New("invalid message format")
	ErrMissingEmail          = errors.New("missing email")
	ErrMissingId             = errors.New("missing id")
	ErrMissingUserName       = errors.New("missing userName")
	ErrNotFoundAuthorization = errors.New("authorization not found ")
	ErrBlockedUser           = errors.New("cannot follow due to block relationship")
	ErrMissingActorID        = errors.New("missing actor_id")
	ErrMissingConversationID = errors.New("missing conversation_id")
	ErrInvalidBlockedPairs   = errors.New("invalid blocked_pairs")
)
