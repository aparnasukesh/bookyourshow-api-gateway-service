package admin

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
	r.POST("/login", h.logIn)
	auth := r.Use(h.authHandler.AdminAuthMiddleware())
	// Theater
	auth.POST("/theater", h.addTheater)
	auth.DELETE("/theater/:id", h.deleteTheaterByID)
	auth.DELETE("/theater", h.deleteTheaterByName)
	auth.GET("/theater/:id", h.getTheaterByID)
	auth.GET("/theater", h.getTheaterByName)
	auth.PUT("/theater/:id", h.updateTheater)
	auth.GET("/theaters", h.listTheaters)
	//Movies
	auth.GET("/movies", h.listMovies)
	// Theater-Types
	auth.GET("/theater/types", h.listTheaterTypes)
	//Screen type
	auth.GET("/screen/types", h.listScreenTypes)
	//Seat categories
	auth.GET("/seat/categories", h.listSeatCategories)

}
func (h *Handler) register(ctx *gin.Context) {
	userData := Admin{}
	if err := ctx.BindJSON(&userData); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	if err := ValidateAdmin(userData); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	if err := h.svc.Register(ctx.Request.Context(), &userData); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "admin registration is pending approval.")
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

// Theater
func (h *Handler) addTheater(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("Authorization")
	userId, err := h.svc.GetUserIDFromToken(ctx, authorization)
	if err != nil {
		h.responseWithError(ctx, http.StatusUnauthorized, errors.New("unauthorized: Invalid token or user ID extraction failed"))
		return
	}
	theater := &Theater{}
	if err := ctx.ShouldBindJSON(&theater); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	theater.OwnerID = uint(userId)
	err = h.svc.AddTheater(ctx, *theater)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "theater added successfully")
}

func (h *Handler) deleteTheaterByID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	err = h.svc.DeleteTheaterByID(ctx, id)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotModified, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "theater deleted successfully")
}

func (h *Handler) deleteTheaterByName(ctx *gin.Context) {
	theaterName := ctx.DefaultQuery("name", "")
	err := h.svc.DeleteTheaterByName(ctx, theaterName)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotModified, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "theater deleted successfully")
}

func (h *Handler) getTheaterByID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	theater, err := h.svc.GetTheaterByID(ctx, id)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get theater details successfully", theater)
}

func (h *Handler) getTheaterByName(ctx *gin.Context) {
	name := ctx.DefaultQuery("name", "")
	theater, err := h.svc.GetTheaterByName(ctx, name)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get theater details successfully", theater)
}

func (h *Handler) updateTheater(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	theater := &Theater{}
	if err := ctx.ShouldBindJSON(&theater); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	err = h.svc.UpdateTheater(ctx, id, *theater)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotModified, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "theater updated successfully")
}

func (h *Handler) listTheaters(ctx *gin.Context) {
	theaters, err := h.svc.ListTheaters(ctx)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNoContent, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list theaters successfully", theaters)
}

// Movies
func (h *Handler) listMovies(ctx *gin.Context) {
	movies, err := h.svc.ListMovies(ctx)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNoContent, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list movies succesfully", movies)
}

// Theater-types
func (h *Handler) listTheaterTypes(ctx *gin.Context) {
	theaterTypes, err := h.svc.ListTheaters(ctx)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNoContent, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list theater-types succesfully", theaterTypes)
}

//Screen types

func (h *Handler) listScreenTypes(ctx *gin.Context) {
	screenTypes, err := h.svc.ListScreenTypes(ctx)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNoContent, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list screen-types successfully", screenTypes)
}

//Seat categories

func (h *Handler) listSeatCategories(ctx *gin.Context) {
	seatCategories, err := h.svc.ListSeatCategories(ctx)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNoContent, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list seat-categories successfully", seatCategories)
}
