package grpcclient

import (
	pb "github.com/aparnasukesh/inter-communication/auth"
	"google.golang.org/grpc"
)

func NewJWT_TokenServiceClient(port string) (pb.JWT_TokenServiceClient, error) {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewJWT_TokenServiceClient(conn), nil
}

func NewUserAuthServiceClient(port string) (pb.UserAuthServiceClient, error) {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewUserAuthServiceClient(conn), nil
}

func NewAdminAuthServiceClient(port string) (pb.AdminAuthServiceClient, error) {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewAdminAuthServiceClient(conn), nil
}

func NewSuperAdminAuthServiceClient(port string) (pb.SuperAdminAuthServiceClient, error) {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewSuperAdminAuthServiceClient(conn), nil
}
