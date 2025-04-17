package messaging

var allowedMessageTypes = map[ServiceType][]MessageType{
	EmailService: {
		"active_user",
		"forgot_password",
		"user_created",
	},
	UserService: {
		"user_created",
		"user_followed",
		"follow_request",
		"unfollow_request",
	},
}

func isAllowedMessageType(service ServiceType, messageType MessageType) bool {
	types, ok := allowedMessageTypes[service]
	if !ok {
		return true
	}
	for _, t := range types {
		if t == messageType {
			return true
		}
	}
	return false
}
