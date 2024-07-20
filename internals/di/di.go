package di

import (
	"github.com/aparnasukesh/api-gateway/config"
	"github.com/aparnasukesh/api-gateway/internals/app/admin"
	"github.com/aparnasukesh/api-gateway/internals/app/middleware"
	superadmin "github.com/aparnasukesh/api-gateway/internals/app/super-admin"
	"github.com/aparnasukesh/api-gateway/internals/app/user"
	grpcclient "github.com/aparnasukesh/api-gateway/pkg/grpcClient"
)

func InitUserModule(cfg config.Config) (*user.Handler, error) {
	pb, err := grpcclient.NewUserGrpcClient(cfg.UserSvcPort)
	if err != nil {
		return nil, err
	}
	auth, err := grpcclient.NewJWT_TokenServiceClient(cfg.AuthSvcPort)
	svc := user.NewService(pb, auth)
	userHandler := user.NewHttpHandler(svc)
	return userHandler, nil
}

func InitAdminModule(cfg config.Config) (*admin.Handler, error) {
	pb, err := grpcclient.NewAdminGrpcClient(cfg.UserSvcPort)
	if err != nil {
		return nil, err
	}
	svc := admin.NewService(pb)
	adminHandler := admin.NewHttpHandler(svc)
	return adminHandler, nil
}

func InitSuperAdminModule(cfg config.Config) (*superadmin.Handler, error) {
	pb, err := grpcclient.NewAdminGrpcClient(cfg.UserSvcPort)
	if err != nil {
		return nil, err
	}
	svc := superadmin.NewService(pb)
	adminHandler := superadmin.NewHttpHandler(svc)
	return adminHandler, nil
}

func InitAuthMiddlewareModule(cfg config.Config) (*middleware.Handler, error) {
	pb, err := grpcclient.NewUserAuthServiceClient(cfg.AuthSvcPort)
	if err != nil {
		return nil, err
	}
	svc := middleware.NewService(pb)
	middlewareHandler := middleware.NewHttpHandler(svc)
	return middlewareHandler, nil
}
