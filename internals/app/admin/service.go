package admin

import (
	"context"

	"github.com/aparnasukesh/inter-communication/user_admin"
)

type Service interface {
	Register(ctx context.Context, signUpData *Admin) error
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

func (s *service) Register(ctx context.Context, signUpData *Admin) error {
	reqData := user_admin.RegisterAdminRequest{
		Username:  signUpData.Username,
		Password:  signUpData.Password,
		Phone:     signUpData.PhoneNumber,
		Email:     signUpData.Email,
		FirstName: signUpData.FirstName,
		LastName:  signUpData.LastName,
		Gender:    signUpData.Gender,
	}
	if _, err := s.pb.RegisterAdmin(ctx, &reqData); err != nil {
		return err
	}
	return nil
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
