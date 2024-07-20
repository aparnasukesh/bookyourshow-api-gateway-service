package superadmin

import (
	"context"

	"github.com/aparnasukesh/inter-communication/user_admin"
)

type Service interface {
	Login(ctx context.Context, loginData *Admin) (string, error)
}

type service struct {
	pb user_admin.AdminServiceClient
}

func NewService(pb user_admin.AdminServiceClient) Service {
	return &service{
		pb: pb,
	}
}
func (s *service) Login(ctx context.Context, loginData *Admin) (string, error) {
	req := user_admin.LoginAdminRequest{
		Email:    loginData.Email,
		Password: loginData.Password,
	}
	res, err := s.pb.LoginAdmin(ctx, &req)
	if err != nil {
		return "", err
	}
	return res.Token, nil
}
