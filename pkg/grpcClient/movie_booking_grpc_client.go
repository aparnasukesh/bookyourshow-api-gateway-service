package grpcclient

import (
	pb "github.com/aparnasukesh/inter-communication/movie_booking"
	"google.golang.org/grpc"
)

func NewMovieBookingGrpcClint(port string) (pb.MovieServiceClient, error) {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewMovieServiceClient(conn), nil
}
