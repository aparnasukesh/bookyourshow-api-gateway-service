package superadmin

import (
	"errors"
	"net/http"

	"github.com/aparnasukesh/api-gateway/pkg/common"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc         Service
	authHandler common.Middleware
}

func NewHttpHandler(svc Service, auth common.Middleware) *Handler {
	return &Handler{
		svc:         svc,
		authHandler: auth,
	}
}

func (h *Handler) MountRoutes(r *gin.RouterGroup) {
	r.POST("/login", h.logIn)

	auth := r.Use(h.authHandler.SuperAdminAuthMiddleware())

	auth.GET("/admin/requests", h.listAdminRequests)
	auth.PUT("/admin/approval", h.adminApproval)
	auth.POST("/movie/register", h.registerMovie)

}
func (h *Handler) logIn(ctx *gin.Context) {
	userData := Admin{}
	if err := ctx.ShouldBindJSON(&userData); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	token, err := h.svc.Login(ctx, &userData)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "login succesfull", token)
}

func (h *Handler) listAdminRequests(ctx *gin.Context) {
	adminLists, err := h.svc.ListAdminRequests(ctx)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return

	}
	h.responseWithData(ctx, http.StatusOK, "retrive admin requests list successfull", adminLists)
}

func (h *Handler) adminApproval(ctx *gin.Context) {
	email := ctx.Query("email")
	is_verified := ctx.Query("is_verified")
	err := h.svc.AdminApproval(ctx, email, is_verified)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return

	}
	h.response(ctx, http.StatusOK, "admin approval successfull")
}

func (h *Handler) registerMovie(ctx *gin.Context) {
	movie := &Movie{}
	if err := ctx.ShouldBindBodyWithJSON(&movie); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	movieId, err := h.svc.RegisterMovie(ctx, *movie)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "movie successfully created", movieId)
}
