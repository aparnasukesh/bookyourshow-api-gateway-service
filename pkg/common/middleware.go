package common

import "github.com/gin-gonic/gin"

type Middleware interface {
	UserAuthMiddleware() gin.HandlerFunc
	AdminAuthMiddleware() gin.HandlerFunc
	SuperAdminAuthMiddleware() gin.HandlerFunc
}
