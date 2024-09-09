package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/aparnasukesh/api-gateway/pkg/common"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc         Service
	authHandler common.Middleware
}

func NewHttpHandler(svc Service, authHandler common.Middleware) *Handler {
	return &Handler{
		svc:         svc,
		authHandler: authHandler,
	}
}

func (h *Handler) MountRoutes(r *gin.RouterGroup) {
	r.POST("/register", h.register)
	r.POST("/register/validate", h.registerValidate)
	r.POST("/login", h.logIn)
	r.POST("/forgot/password", h.forgotPassword)
	r.POST("/reset/password", h.resetPassword)

	auth := r.Use(h.authHandler.UserAuthMiddleware())
	auth.GET("/profile", h.getProfile)
	auth.PUT("/profile/:id", h.updateUserProfile)

}

func (h *Handler) register(ctx *gin.Context) {
	userData := User{}
	if err := ctx.BindJSON(&userData); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	if err := ValidateUser(userData); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	if err := h.svc.Register(ctx.Request.Context(), &userData); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "signup successfull", map[string]string{
		"redirect": "http://localhost:8080/gateway/user/register/validate",
	})
}

func (h *Handler) registerValidate(ctx *gin.Context) {
	userData := User{}
	if err := ctx.ShouldBindJSON(&userData); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	if err := h.svc.RegisterValidate(ctx, &userData); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "register validate successfull")
}

func (h *Handler) logIn(ctx *gin.Context) {
	userData := User{}
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

func (h *Handler) getProfile(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("Authorization")
	userId, err := h.svc.GetUserIDFromToken(ctx, authorization)
	if err != nil {
		h.responseWithError(ctx, http.StatusUnauthorized, errors.New("unauthorized: Invalid token or user ID extraction failed"))
		return
	}
	profileDetails, err := h.svc.GetProfile(ctx, userId)
	if err != nil {
		h.responseWithError(ctx, http.StatusNotFound, errors.New("profile not found: Unable to retrieve profile details for the user"))
		return
	}

	h.responseWithData(ctx, http.StatusOK, "User profile details retrieved successfully", profileDetails)
}

func (h *Handler) updateUserProfile(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	user := &UserProfileDetails{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	err = h.svc.UpdateUserProfile(ctx, id, *user)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "update user profile successfull")
}

func (h *Handler) forgotPassword(ctx *gin.Context) {
	email := ForgotPassword{}
	if err := ctx.ShouldBindJSON(&email); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	err := h.svc.ForgotPassword(ctx, email.Email)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "otp send successfull")
}

func (h *Handler) resetPassword(ctx *gin.Context) {
	data := ResetPassword{}
	if err := ctx.ShouldBindJSON(&data); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	err := h.svc.ResetPassword(ctx, data)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "password reset successfull")
}
