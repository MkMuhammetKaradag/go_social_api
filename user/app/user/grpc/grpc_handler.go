package usergrpc

import (
	"context"
	"socialmedia/user/proto/userpb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
}

// GetUser implements the gRPC method from the proto file.
func (s *UserService) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	// Simülasyon (örnek veri)
	return &userpb.GetUserResponse{
		Id:        req.UserId,
		Username:  "john_doe",
		AvatarUrl: wrapperspb.String("https://example.com/avatar.png"),
	}, nil
}
