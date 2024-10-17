package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/aparnasukesh/api-gateway/pkg/common"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHttpHandler(svc Service) common.Middleware {
	return &Handler{
		svc: svc,
	}
}
func (h *Handler) UserAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.Request.Header.Get("Authorization")
		if authorization == "" {
			h.responseWithError(ctx, http.StatusUnauthorized, fmt.Errorf("authorization header is missing"))
			ctx.Abort()
			return
		}

		err := h.svc.UserAuthentication(ctx, authorization)
		if err != nil {
			formattedError := ExtractErrorMessage(err)
			h.responseWithError(ctx, http.StatusUnauthorized, errors.New(formattedError))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
func (h *Handler) AdminAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.Request.Header.Get("Authorization")
		if err := h.svc.AdminAuthentication(ctx, authorization); err != nil {
			formattedError := ExtractErrorMessage(err)
			h.responseWithError(ctx, http.StatusUnauthorized, errors.New(formattedError))
			ctx.Abort()
			return
		}
		ctx.Next()

	}
}

func (h *Handler) SuperAdminAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.Request.Header.Get("Authorization")
		if err := h.svc.SuperAdminAuthentication(ctx, authorization); err != nil {
			formattedError := ExtractErrorMessage(err)
			h.responseWithError(ctx, http.StatusUnauthorized, errors.New(formattedError))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func (h *Handler) UserPaymentAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		token, err := ctx.Cookie("UserAuthorization")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Success": false,
				"Message": "Authorization failed",
				"Error":   err.Error(),
			})
			ctx.Abort()
			return
		}
		if token == "" {
			h.responseWithError(ctx, http.StatusUnauthorized, fmt.Errorf("authorization header is missing"))
			ctx.Abort()
			return
		}

		err = h.svc.UserAuthentication(ctx, token)
		if err != nil {
			formattedError := ExtractErrorMessage(err)
			h.responseWithError(ctx, http.StatusUnauthorized, errors.New(formattedError))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
