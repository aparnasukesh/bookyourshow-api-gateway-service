package boot

import (
	"log"

	"github.com/aparnasukesh/api-gateway/config"
	"github.com/aparnasukesh/api-gateway/internals/di"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type resources struct {
	cfg config.Config
}

func Start(cfg config.Config) {
	res := &resources{cfg: cfg}
	r := gin.Default()

	res.MountRoutes(r)

	r.Run(":8080")
}

func (m resources) MountRoutes(r *gin.Engine) {
	userHandler, err := di.InitUserModule(m.cfg)
	if err != nil {
		log.Fatalf("Error happened while user module initialization: %v", err)
	}
	adminHandler, err := di.InitAdminModule(m.cfg)
	if err != nil {
		log.Fatalf("Error happened while admin module initialization: %v", err)
	}
	superAdminHandler, err := di.InitSuperAdminModule(m.cfg)
	if err != nil {
		log.Fatalf("Error happend while super admin module initialization: %v", err)
	}

	r.Use(cors.New(SetCors()))

	gateway := r.Group("/gateway")
	{
		user := gateway.Group("/user")
		userHandler.MountRoutes(user)

		admin := gateway.Group("/admin")
		adminHandler.MountRoutes(admin)

		superAdmin := gateway.Group("/superadmin")
		superAdminHandler.MountRoutes(superAdmin)
	}

}

func SetCors() cors.Config {
	return cors.Config{
		AllowOrigins:     []string{"https://api.bookyourshow.com", "*"}, // Replace with actual Razorpay URL or use "*" to allow all
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE", "OPTION"},
	}
}
