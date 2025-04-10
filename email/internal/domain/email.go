package domain

import "socialmedia/shared/messaging"

type EmailType string

const (
	EmailTypeActivation     EmailType = "active_user"
	EmailTypeForgotPassword EmailType = "forgot_password"
)

type EmailRequest struct {
	Type messaging.MessageType `json:"type"`
	Data interface{}           `json:"data"`
}

type Mailer interface {
	Send(to, subject, template string) error
}
