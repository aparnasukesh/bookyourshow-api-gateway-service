package superadmin

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

func (h *Handler) MountRout(r *gin.RouterGroup) {
	r.POST("/login", h.logIn)

}
func (h *Handler) logIn(ctx *gin.Context) {
	userData := Admin{}
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
