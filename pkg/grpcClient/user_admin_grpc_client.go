package grpcclient

import (
	pb "github.com/aparnasukesh/inter-communication/user_admin"
	"google.golang.org/grpc"
)

func NewUserGrpcClient(port string) (pb.UserServiceClient, error) {
	conn, err := grpc.Dial("user-admin-svc:"+port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewUserServiceClient(conn), nil
}

func NewAdminGrpcClient(port string) (pb.AdminServiceClient, error) {
	conn, err := grpc.Dial("user-admin-svc:"+port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewAdminServiceClient(conn), nil
}

func NewSuperAdminServiceClient(port string) (pb.SuperAdminServiceClient, error) {
	conn, err := grpc.Dial("user-admin-svc:"+port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewSuperAdminServiceClient(conn), nil
}
