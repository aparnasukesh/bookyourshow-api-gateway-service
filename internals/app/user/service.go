package user

import (
	"context"
	"strings"

	"github.com/aparnasukesh/inter-communication/auth"
	"github.com/aparnasukesh/inter-communication/user_admin"
)

type Service interface {
	Register(ctx context.Context, signUpData *User) error
	RegisterValidate(ctx context.Context, userData *User) error
	Login(ctx context.Context, loginData *User) (string, error)
	GetUserIDFromToken(ctx context.Context, authorization string) (int, error)
	GetProfile(ctx context.Context, userId int) (*UserProfileDetails, error)
}

type service struct {
	pb   user_admin.UserServiceClient
	auth auth.JWT_TokenServiceClient
}

func NewService(pb user_admin.UserServiceClient, auth auth.JWT_TokenServiceClient) Service {
	return &service{
		pb:   pb,
		auth: auth,
	}
}

func (s *service) Register(ctx context.Context, signUpData *User) error {
	reqData := user_admin.RegisterUserRequest{
		Username:  signUpData.Username,
		Password:  signUpData.Password,
		Phone:     signUpData.PhoneNumber,
		Email:     signUpData.Email,
		FirstName: signUpData.FirstName,
		LastName:  signUpData.LastName,
		Gender:    signUpData.Gender,
	}
	if _, err := s.pb.RegisterUser(ctx, &reqData); err != nil {
		return err
	}
	return nil
}

func (s *service) RegisterValidate(ctx context.Context, userData *User) error {
	req := user_admin.ValidateUserRequest{
		Email: userData.Email,
		Otp:   userData.Otp,
	}
	if _, err := s.pb.ValidateUser(ctx, &req); err != nil {
		return err
	}
	return nil
}
func (s *service) Login(ctx context.Context, loginData *User) (string, error) {
	req := user_admin.LoginUserRequest{
		Email:    loginData.Email,
		Password: loginData.Password,
	}
	res, err := s.pb.LoginUser(ctx, &req)
	if err != nil {
		return "", err
	}
	return res.Token, nil
}

func (s *service) GetUserIDFromToken(ctx context.Context, authorization string) (int, error) {
	tokenParts := strings.Split(authorization, "Bearer ")
	token := tokenParts[1]
	var userId int

	response, err := s.auth.GetUserID(ctx, &auth.GetUserIDRequest{
		Token: token,
	})
	userId = int(response.UserId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (s *service) GetProfile(ctx context.Context, userId int) (*UserProfileDetails, error) {
	response, err := s.pb.GetUserProfile(ctx, &user_admin.GetProfileRequest{
		UserId: int32(userId),
	})
	if err != nil {
		return nil, err
	}
	details, err := BuildGetUserProfile(response)
	if err != nil {
		return nil, err
	}
	return details, nil
}
