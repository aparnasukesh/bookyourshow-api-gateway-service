package superadmin

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
	auth.PUT("/movie/:id", h.updateMovie)
	auth.GET("/movies", h.listMovies)
	auth.GET("/movie/:id", h.getMovieDetails)
	auth.DELETE("/movie/:id", h.deleteMovie)

	auth.POST("/theater/type", h.addTheaterType)
	auth.DELETE("/theater/type/:id", h.deleteTheaterTypeById)
	auth.DELETE("/theater/type", h.deleteTheaterTypeByName)
	auth.GET("/theater/type/:id", h.getTheaterTypeByID)
	auth.GET("/theater/type", h.getTheaterTypeByName)
	auth.PUT("/theater/type/:id", h.updateTheaterType)
	auth.GET("/theater/types", h.listTheaterTypes)
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
	approval := &AdminApproval{}
	if err := ctx.ShouldBindJSON(&approval); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return

	}
	err := h.svc.AdminApproval(ctx, approval.Email, approval.IsVerified)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return

	}
	h.response(ctx, http.StatusOK, "admin approval successfull")
}

// movies
func (h *Handler) registerMovie(ctx *gin.Context) {
	movie := &Movie{}
	if err := ctx.ShouldBindJSON(&movie); err != nil {
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

func (h *Handler) updateMovie(ctx *gin.Context) {
	movie := &Movie{}
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}

	if err := ctx.ShouldBindJSON(&movie); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	if err := h.svc.UpdateMovie(ctx, *movie, id); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotModified, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "movie updated succesfully")
}

func (h *Handler) listMovies(ctx *gin.Context) {
	movies, err := h.svc.ListMovies(ctx)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNoContent, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list movies succesfully", movies)
}

func (h *Handler) getMovieDetails(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	movie, err := h.svc.GetMovieDetails(ctx, id)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotModified, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get movie details succesfully", movie)
}

func (h *Handler) deleteMovie(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	err = h.svc.DeleteMovie(ctx, id)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotModified, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "movie deleted succesfully")
}

// theater-type
func (h *Handler) addTheaterType(ctx *gin.Context) {
	theaterType := &TheaterType{}
	if err := ctx.ShouldBindJSON(&theaterType); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	err := h.svc.AddTheaterType(ctx, *theaterType)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "theater type addedd successfully")
}

func (h *Handler) deleteTheaterTypeById(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	err = h.svc.DeleteTheaterTypeById(ctx, id)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotModified, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "theater type  deleted succesfully")
}

func (h *Handler) deleteTheaterTypeByName(ctx *gin.Context) {
	theaterName := ctx.DefaultQuery("name", "")
	err := h.svc.DeleteTheaterTypeByName(ctx, theaterName)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotModified, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "theater type  deleted succesfully")
}

func (h *Handler) getTheaterTypeByID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	theaterType, err := h.svc.GetTheaterTypeByID(ctx, id)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get theater-type details succesfully", theaterType)
}

func (h *Handler) getTheaterTypeByName(ctx *gin.Context) {
	name := ctx.DefaultQuery("name", "")
	theaterType, err := h.svc.GetTheaterTypeByName(ctx, name)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get theater-type details succesfully", theaterType)
}

func (h *Handler) updateTheaterType(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	theatertype := &TheaterType{}
	if err := ctx.ShouldBindJSON(&theatertype); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	err = h.svc.UpdateTheaterType(ctx, id, *theatertype)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotModified, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "get theater-type details succesfully")
}

func (h *Handler) listTheaterTypes(ctx *gin.Context) {
	theaterTypes, err := h.svc.ListMovies(ctx)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNoContent, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list theater-types succesfully", theaterTypes)
}
