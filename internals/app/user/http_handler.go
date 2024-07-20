package user

import (
	"errors"
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

func (h *Handler) MountRoutes(r *gin.RouterGroup) {
	r.POST("/register", h.register)
	r.POST("/register/validate", h.registerValidate)
	r.POST("/login", h.logIn)

	authorized := r.Group("/")
	authorized.GET("/profile", h.getProfile)
}

func (h *Handler) register(ctx *gin.Context) {
	userData := User{}
	if err := ctx.BindJSON(&userData); err != nil {
		h.responseWithError(ctx, http.StatusBadRequest, err)
		return
	}
	if err := ValidateUser(userData); err != nil {
		h.responseWithError(ctx, http.StatusBadRequest, err)
		return
	}
	if err := h.svc.Register(ctx.Request.Context(), &userData); err != nil {
		h.responseWithError(ctx, http.StatusNotFound, err)
		return
	}
	h.responseWithData(ctx, http.StatusOK, "signup successfull", map[string]string{
		"redirect": "http://localhost:8080/gateway/user/register/validate",
	})
}

func (h *Handler) registerValidate(ctx *gin.Context) {
	userData := User{}
	if err := ctx.ShouldBindJSON(&userData); err != nil {
		h.responseWithError(ctx, http.StatusBadRequest, err)
		return
	}
	if err := h.svc.RegisterValidate(ctx, &userData); err != nil {
		h.responseWithError(ctx, http.StatusNotFound, err)
		return
	}
	h.response(ctx, http.StatusOK, "register validate successfull")
}

func (h *Handler) logIn(ctx *gin.Context) {
	userData := User{}
	if err := ctx.ShouldBindJSON(&userData); err != nil {
		h.responseWithError(ctx, http.StatusBadRequest, err)
		return
	}
	token, err := h.svc.Login(ctx, &userData)
	if err != nil {
		h.responseWithError(ctx, http.StatusNotFound, err)
		return
	}
	h.responseWithData(ctx, http.StatusOK, "login succesfull", token)
}

func (h *Handler) getProfile(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("Authorization")
	userId, err := h.svc.GetUserIDFromToken(ctx, authorization)
	if err != nil {
		h.responseWithError(ctx, http.StatusUnauthorized, errors.New("Unauthorized: Invalid token or user ID extraction failed"))
		return
	}
	profileDetails, err := h.svc.GetProfile(ctx, userId)
	if err != nil {
		h.responseWithError(ctx, http.StatusNotFound, errors.New("Profile not found: Unable to retrieve profile details for the user"))
		return
	}

	h.responseWithData(ctx, http.StatusOK, "User profile details retrieved successfully", profileDetails)

}
