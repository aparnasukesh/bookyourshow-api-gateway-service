package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHttpHandler(svc Service) *Handler {
	return &Handler{
		svc: svc,
	}
}
func (h *Handler) UserAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.Request.Header.Get("Authorization")
		if err := h.svc.UserAuthentication(ctx, authorization); err != nil {
			h.responseWithError(ctx, http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
