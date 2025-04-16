package messaging

type ServiceType string
type MessageType string

const (
	AuthService   ServiceType = "auth"
	UserService   ServiceType = "user"
	FollowService ServiceType = "follow"
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
	UserCreated          MessageType
	UserFollowed         MessageType
	FollowRequestCreated MessageType
}{
	UserCreated:          "user_created",
	UserFollowed:         "user_followed",
	FollowRequestCreated: "follow_request_created",
}

var FollowTypes = struct {
	UserCreated MessageType
}{
	UserCreated: "user_created",
}
