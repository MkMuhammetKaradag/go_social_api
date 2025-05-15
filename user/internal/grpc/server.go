package grpcserver

import (
	"log"
	"net"
	usergrpc "socialmedia/user/app/user/grpc"
	"socialmedia/user/proto/userpb"

	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
}

func NewGRPCServer() *Server {
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &usergrpc.UserService{}) // implement edilmiş struct

	return &Server{
		grpcServer: grpcServer,
	}
}

func (s *Server) Start(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server listening at %v", address)
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
