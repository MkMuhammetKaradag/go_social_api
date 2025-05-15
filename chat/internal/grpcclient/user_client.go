package grpcclient

import (
	"context"
	"fmt"
	"socialmedia/chat/proto/userpb" // proto'ları buraya göre import et

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	conn   *grpc.ClientConn
	Client userpb.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	client := userpb.NewUserServiceClient(conn)
	return &UserClient{
		conn:   conn,
		Client: client,
	}, nil
}

func (u *UserClient) Close() error {
	return u.conn.Close()
}

// Örnek bir çağrı fonksiyonu
func (u *UserClient) GetUserByID(ctx context.Context, userID string) (*userpb.GetUserResponse, error) {
	return u.Client.GetUser(ctx, &userpb.GetUserRequest{UserId: userID})
}
