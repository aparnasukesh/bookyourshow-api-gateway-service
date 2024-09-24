package grpcclient

import (
	pb "github.com/aparnasukesh/inter-communication/movie_booking"
	"google.golang.org/grpc"
)

func NewMovieBookingGrpcClint(port string) (pb.MovieServiceClient, pb.TheatreServiceClient, error) {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return pb.NewMovieServiceClient(conn), pb.NewTheatreServiceClient(conn), nil
}
