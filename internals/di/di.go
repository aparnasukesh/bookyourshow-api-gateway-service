package di

import (
	"log"

	"github.com/aparnasukesh/api-gateway/config"
	"github.com/aparnasukesh/api-gateway/internals/app/admin"
	"github.com/aparnasukesh/api-gateway/internals/app/middleware"
	superadmin "github.com/aparnasukesh/api-gateway/internals/app/super-admin"
	"github.com/aparnasukesh/api-gateway/internals/app/user"
	"github.com/aparnasukesh/api-gateway/pkg/common"
	grpcclient "github.com/aparnasukesh/api-gateway/pkg/grpcClient"
	"github.com/aparnasukesh/api-gateway/pkg/rabbitmq"
)

func InitUserModule(cfg config.Config) (*user.Handler, error) {
	pb, err := grpcclient.NewUserGrpcClient(cfg.UserSvcPort)
	if err != nil {
		return nil, err
	}
	authHandler, err := InitAuthMiddlewareModule(cfg)
	if err != nil {
		log.Fatalf("Error happened while authmiddleware module initialization")
	}

	auth, err := grpcclient.NewJWT_TokenServiceClient(cfg.AuthSvcPort)
	if err != nil {
		log.Fatalf("Error happened while TokenServiceClient module initialization")
	}
	movieBooking, theater, booking, err := grpcclient.NewMovieBookingGrpcClint(cfg.MovieBookingPort)

	paymentClient, err := grpcclient.NewBookingPaymentServiceClient(cfg.PaymentPort)
	rabbitmqConnection, err := rabbitmq.NewRabbitMQConnection()
	if err != nil {
		return nil, err
	}
	svc := user.NewService(pb, auth, movieBooking, theater, booking, paymentClient, rabbitmqConnection)
	userHandler := user.NewHttpHandler(svc, authHandler)
	return userHandler, nil
}

func InitAdminModule(cfg config.Config) (*admin.Handler, error) {
	pb, err := grpcclient.NewAdminGrpcClient(cfg.UserSvcPort)
	if err != nil {
		return nil, err
	}
	authHandler, err := InitAuthMiddlewareModule(cfg)
	if err != nil {
		log.Fatalf("Error happpened while authmiddleware module initialization")
	}
	auth, err := grpcclient.NewJWT_TokenServiceClient(cfg.AuthSvcPort)
	if err != nil {
		log.Fatalf("Error happened while TokenServiceClient module initialization")
	}
	svc := admin.NewService(pb, auth)
	adminHandler := admin.NewHttpHandler(svc, authHandler)
	return adminHandler, nil
}

func InitSuperAdminModule(cfg config.Config) (*superadmin.Handler, error) {
	pb, err := grpcclient.NewSuperAdminServiceClient(cfg.UserSvcPort)
	if err != nil {
		return nil, err
	}
	authHandler, err := InitAuthMiddlewareModule(cfg)
	if err != nil {
		log.Fatalf("Error happpened while authmiddleware module initialization")
	}
	auth, err := grpcclient.NewJWT_TokenServiceClient(cfg.AuthSvcPort)
	if err != nil {
		log.Fatalf("Error happened while TokenServiceClient module initialization")
	}
	movieBooking, _, _, err := grpcclient.NewMovieBookingGrpcClint(cfg.MovieBookingPort)
	svc := superadmin.NewService(pb, auth, movieBooking)
	adminHandler := superadmin.NewHttpHandler(svc, authHandler)
	return adminHandler, nil
}

func InitAuthMiddlewareModule(cfg config.Config) (common.Middleware, error) {
	userSvcClient, err := grpcclient.NewUserAuthServiceClient(cfg.AuthSvcPort)
	if err != nil {
		return nil, err
	}
	adminSvcClient, err := grpcclient.NewAdminAuthServiceClient(cfg.AuthSvcPort)
	if err != nil {
		return nil, err
	}
	superAdminSvcClient, err := grpcclient.NewSuperAdminAuthServiceClient(cfg.AuthSvcPort)
	if err != nil {
		return nil, err
	}
	svc := middleware.NewService(userSvcClient, adminSvcClient, superAdminSvcClient)
	middlewareHandler := middleware.NewHttpHandler(svc)
	return middlewareHandler, nil
}
