package messaging

type ServiceType string
type MessageType string

const (
	AuthService   ServiceType = "auth"
	UserService   ServiceType = "user"
	FallowService ServiceType = "fallow"
	EmailService  ServiceType = "email"
	ChatService   ServiceType = "chat"
)

var EmailTypes = struct {
	ActivateUser   MessageType
	ForgotPassword MessageType
}{
	ActivateUser:   "active_user",
	ForgotPassword: "forgot_password",
}

var UserTypes = struct {
	UserCreated MessageType
}{
	UserCreated: "user_created",
}

var FallowTypes = struct {
	UserCreated MessageType
}{
	UserCreated: "user_created",
}
