package middleware

import (
	"context"

	pb "github.com/aparnasukesh/inter-communication/auth"
)

type Service interface {
	UserAuthentication(ctx context.Context, token string) error
	AdminAuthentication(ctx context.Context, token string) error
	SuperAdminAuthentication(ctx context.Context, token string) error
}

type service struct {
	userSvc       pb.UserAuthServiceClient
	adminSvc      pb.AdminAuthServiceClient
	superAdminSvc pb.SuperAdminAuthServiceClient
}

func NewService(userSvc pb.UserAuthServiceClient, adminSvc pb.AdminAuthServiceClient, superAdminSvc pb.SuperAdminAuthServiceClient) Service {
	return &service{
		userSvc:       userSvc,
		adminSvc:      adminSvc,
		superAdminSvc: superAdminSvc,
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

func (s *service) AdminAuthentication(ctx context.Context, token string) error {
	if _, err := s.adminSvc.AdminAuthRequired(ctx, &pb.AuthRequest{
		Token: token,
	}); err != nil {
		return err
	}
	return nil
}

func (s *service) SuperAdminAuthentication(ctx context.Context, token string) error {
	if _, err := s.superAdminSvc.SuperAdminAuthRequired(ctx, &pb.AuthRequest{
		Token: token,
	}); err != nil {
		return err
	}
	return nil
}
