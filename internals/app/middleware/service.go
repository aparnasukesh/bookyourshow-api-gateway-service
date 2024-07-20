package middleware

import (
	"context"

	pb "github.com/aparnasukesh/inter-communication/auth"
)

type Service interface {
	UserAuthentication(ctx context.Context, token string) error
	// AdminAuthentication(ctx context.Context, token string) error
	// SuperAdminAuthentication(ctx context.Context, token string) error
}

type service struct {
	userSvc pb.UserAuthServiceClient
	// adminSvc      pb.AdminAuthServiceClient
	// superAdminSvc pb.SuperAdminServiceClient
}

func NewService(userSvc pb.UserAuthServiceClient) Service {
	return &service{
		userSvc: userSvc,
	}
}

func (s *service) UserAuthentication(ctx context.Context, token string) error {
	if _, err := s.userSvc.UserAuthRequired(ctx, &pb.AuthRequest{
		Token: token,
	}); err != nil {
		return err
	}
	return nil
}
