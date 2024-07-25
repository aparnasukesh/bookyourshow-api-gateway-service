package superadmin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aparnasukesh/inter-communication/auth"
	"github.com/aparnasukesh/inter-communication/movie_booking"
	"github.com/aparnasukesh/inter-communication/user_admin"
)

type Service interface {
	Login(ctx context.Context, loginData *Admin) (string, error)
	ListAdminRequests(ctx context.Context) ([]AdminRequestResponse, error)
	AdminApproval(ctx context.Context, email string, is_verified string) error
	RegisterMovie(ctx context.Context, movie Movie) (int, error)
}

type service struct {
	pb           user_admin.SuperAdminServiceClient
	movieBooking movie_booking.MovieServiceClient
}

func NewService(pb user_admin.SuperAdminServiceClient, auth auth.JWT_TokenServiceClient, movieBooking movie_booking.MovieServiceClient) Service {
	return &service{
		pb:           pb,
		movieBooking: movieBooking,
	}
}
func (s *service) Login(ctx context.Context, loginData *Admin) (string, error) {
	req := user_admin.LoginSuperAdminRequest{
		Email:    loginData.Email,
		Password: loginData.Password,
	}
	res, err := s.pb.LoginSuperAdmin(ctx, &req)
	if err != nil {
		return "", err
	}
	return res.Token, nil
}
func (s *service) ListAdminRequests(ctx context.Context) ([]AdminRequestResponse, error) {
	res, err := s.pb.ListAdminRequests(ctx, &user_admin.ListAdminRequestsRequest{})
	if err != nil {
		return nil, err
	}

	adminRequests := make([]AdminRequestResponse, len(res.Email))
	for i, admin := range res.Email {
		adminRequests[i] = AdminRequestResponse{
			Email: admin.Email,
		}
	}
	return adminRequests, nil
}

func (s *service) AdminApproval(ctx context.Context, email string, is_verified string) error {
	isVerified, err := strconv.ParseBool(is_verified)
	if err != nil {
		return fmt.Errorf("invalid value for is_verified: %v", err)
	}
	_, err = s.pb.AdminApproval(ctx, &user_admin.AdminApprovalRequest{
		Email:      email,
		IsVerified: isVerified,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) RegisterMovie(ctx context.Context, movie Movie) (int, error) {
	response, err := s.movieBooking.RegisterMovie(ctx, &movie_booking.RegisterMovieRequest{
		Title:       movie.Title,
		Description: movie.Description,
		Duration:    int32(movie.Duration),
		Genre:       movie.Genre,
		ReleaseDate: movie.ReleaseDate,
		Rating:      float32(movie.Rating),
		Language:    movie.Language,
	})
	if err != nil {
		return 0, err
	}
	return int(response.MovieId), nil
}
