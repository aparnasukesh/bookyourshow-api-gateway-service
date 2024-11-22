package grpcclient

import (
	"log"

	pb "github.com/aparnasukesh/inter-communication/auth"
	"google.golang.org/grpc"
)

// func NewJWT_TokenServiceClient(port string) (pb.JWT_TokenServiceClient, error) {
// 	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return pb.NewJWT_TokenServiceClient(conn), nil
// }

func NewJWT_TokenServiceClient(port string) (pb.JWT_TokenServiceClient, error) {
	address := "auth-svc.default.svc.cluster.local:" + port
	serviceConfig := `{"loadBalancingPolicy": "round_robin"}`
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(serviceConfig))
	if err != nil {
		log.Printf("Failed to connect to gRPC service: %v", err)
		return nil, err
	}
	return pb.NewJWT_TokenServiceClient(conn), nil
}

func NewUserAuthServiceClient(port string) (pb.UserAuthServiceClient, error) {
	address := "auth-svc.default.svc.cluster.local:" + port
	serviceConfig := `{"loadBalancingPolicy": "round_robin"}`
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(serviceConfig))
	if err != nil {
		log.Printf("Failed to connect to gRPC service: %v", err)
		return nil, err
	}
	return pb.NewUserAuthServiceClient(conn), nil
}

func NewAdminAuthServiceClient(port string) (pb.AdminAuthServiceClient, error) {
	address := "auth-svc.default.svc.cluster.local:" + port
	serviceConfig := `{"loadBalancingPolicy": "round_robin"}`
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(serviceConfig))
	if err != nil {
		log.Printf("Failed to connect to gRPC service: %v", err)
		return nil, err
	}
	return pb.NewAdminAuthServiceClient(conn), nil
}

func NewSuperAdminAuthServiceClient(port string) (pb.SuperAdminAuthServiceClient, error) {
	address := "auth-svc.default.svc.cluster.local:" + port
	serviceConfig := `{"loadBalancingPolicy": "round_robin"}`
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(serviceConfig))
	if err != nil {
		log.Printf("Failed to connect to gRPC service: %v", err)
		return nil, err
	}
	return pb.NewSuperAdminAuthServiceClient(conn), nil
}
