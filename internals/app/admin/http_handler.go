package admin

import (
	"fmt"
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

func (h *Handler) MountRout(r *gin.RouterGroup) {
	r.POST("/register", h.register)
	r.POST("/login", h.logIn)

}
func (h *Handler) register(ctx *gin.Context) {
	userData := Admin{}
	if err := ctx.BindJSON(&userData); err != nil {
		h.responseWithError(ctx, http.StatusBadRequest, err)
		return
	}
	if err := ValidateAdmin(userData); err != nil {
		h.responseWithError(ctx, http.StatusBadRequest, err)
		return
	}
	if err := h.svc.Register(ctx.Request.Context(), &userData); err != nil {
		h.responseWithError(ctx, http.StatusNotFound, err)
		return
	}
	h.response(ctx, http.StatusOK, "admin registration is pending approval.")
}

func (h *Handler) logIn(ctx *gin.Context) {
	userData := Admin{}
	if err := ctx.ShouldBindJSON(&userData); err != nil {
		h.responseWithError(ctx, http.StatusBadRequest, err)
		return
	}
	token, err := h.svc.Login(ctx, &userData)
	if err != nil {
		fmt.Println("err:", err)
		h.responseWithError(ctx, http.StatusNotFound, err)
		return
	}
	h.responseWithData(ctx, http.StatusOK, "login succesfull", token)
}
