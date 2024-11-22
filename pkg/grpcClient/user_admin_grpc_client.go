package grpcclient

import (
	"log"

	pb "github.com/aparnasukesh/inter-communication/user_admin"
	"google.golang.org/grpc"
)

// func NewUserGrpcClient(port string) (pb.UserServiceClient, error) {
// 	conn, err := grpc.Dial("user-admin-svc:"+port, grpc.WithInsecure(),grpc.WithDefaultServiceConfig({"loadBalancingPolicy":"round_robin"}))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return pb.NewUserServiceClient(conn), nil
// }

func NewUserGrpcClient(port string) (pb.UserServiceClient, error) {
	address := "user-admin-svc.svc.cluster.local:" + port
	serviceConfig := `{"loadBalancingPolicy": "round_robin"}`
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(serviceConfig))
	if err != nil {
		log.Printf("Failed to connect to gRPC service: %v", err)
		return nil, err
	}
	return pb.NewUserServiceClient(conn), nil
}

func NewAdminGrpcClient(port string) (pb.AdminServiceClient, error) {
	address := "user-admin-svc.svc.cluster.local:" + port
	serviceConfig := `{"loadBalancingPolicy": "round_robin"}`
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(serviceConfig))
	if err != nil {
		log.Printf("Failed to connect to gRPC service: %v", err)
		return nil, err
	}
	return pb.NewAdminServiceClient(conn), nil
}

func NewSuperAdminServiceClient(port string) (pb.SuperAdminServiceClient, error) {
	address := "user-admin-svc.svc.cluster.local:" + port
	serviceConfig := `{"loadBalancingPolicy": "round_robin"}`
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(serviceConfig))
	if err != nil {
		log.Printf("Failed to connect to gRPC service: %v", err)
		return nil, err
	}
	return pb.NewSuperAdminServiceClient(conn), nil
}
