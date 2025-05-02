package messaging

type ServiceType string
type MessageType string

const (
	AuthService         ServiceType = "auth"
	UserService         ServiceType = "user"
	FollowService       ServiceType = "follow"
	EmailService        ServiceType = "email"
	ChatService         ServiceType = "chat"
	NotificationService ServiceType = "notification"
)

var EmailTypes = struct {
	ActivateUser   MessageType
	ForgotPassword MessageType
}{
	ActivateUser:   "active_user",
	ForgotPassword: "forgot_password",
}

var UserTypes = struct {
	UserCreated     MessageType
	UserFollowed    MessageType
	FollowRequest   MessageType
	UnFollowRequest MessageType
	UserBlocked     MessageType
	UserUnBlocked   MessageType
	UserUpdated     MessageType
}{
	UserCreated:     "user_created",
	UserFollowed:    "user_followed",
	FollowRequest:   "follow_request",
	UnFollowRequest: "unfollow_request",
	UserBlocked:     "user_blocked",
	UserUnBlocked:   "user_unblocked",
	UserUpdated:     "user_updated",
}

var FollowTypes = struct {
	UserCreated MessageType
}{
	UserCreated: "user_created",
}
var ChatTypes = struct {
	UserBlockedInGroupConversation MessageType
}{
	UserBlockedInGroupConversation: "user_blockedIn_group_conversation",
}
