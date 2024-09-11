package admin

import (
	"errors"
	"net/http"
	"strconv"
	"time"

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
	r.POST("/forgot/password", h.forgotPassword)
	r.POST("/reset/password", h.resetPassword)
	auth := r.Use(h.authHandler.AdminAuthMiddleware())

	auth.GET("/profile/:id", h.getAdminProfile)
	auth.PUT("/profile/:id", h.updateAdminProfile)
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
	//Theater screen
	auth.POST("/theater/screen", h.addTheaterScreen)
	auth.DELETE("/theater/screen/:id", h.deleteTheaterScreenByID)
	auth.DELETE("/theater/screen", h.deleteTheaterScreenByNumber)
	auth.GET("/theater/screen/:id", h.getTheaterScreenByID)
	auth.GET("/theater/screen", h.getTheaterScreenByNumber)
	auth.PUT("/theater/screen/:id", h.updateTheaterScreen)
	auth.GET("/theater/screens", h.listTheaterScreens)
	//Show times
	auth.POST("/showtime", h.addShowtime)
	auth.DELETE("/showtime/:id", h.deleteShowtimeByID)
	auth.DELETE("/showtime", h.deleteShowtimeByDetails)
	auth.GET("/showtime/:id", h.getShowtimeByID)
	auth.GET("/showtime", h.getShowtimeByDetails)
	auth.PUT("/showtime/:id", h.updateShowtime)
	auth.GET("/showtimes", h.listShowtimes)
	// Movie Schedule
	auth.POST("/movie/schedule", h.addMovieSchedule)
	auth.PUT("/movie/schedule/:id", h.updateMovieSchedule)
	auth.GET("/movie/schedules", h.getAllMovieSchedules)
	auth.GET("/movie/schedule/movieid", h.getMovieScheduleByMovieID)
	auth.GET("/movie/schedule/theaterid", h.getMovieScheduleByTheaterID)
	auth.GET("/movie/schedule/movieid/theaterid", h.getMovieScheduleByMovieIdAndTheaterId)
	auth.GET("/movie/schedule/movieid/showtimeid", h.getMovieScheduleByMovieIdAndShowTimeId)
	auth.GET("/movie/schedule/theaterid/showtimeid", h.getMovieScheduleByTheaterIdAndShowTimeId)
	auth.GET("/movie/schedule/:id", h.getMovieScheduleByID)
	auth.DELETE("/movie/schedule/:id", h.deleteMovieScheduleById)
	auth.DELETE("/movie/schedule/movieid/theaterid", h.deleteMovieScheduleByMovieIdAndTheaterId)
	auth.DELETE("/movie/schedule/movieid/theaterid/showtimeid", h.deleteMovieScheduleByMovieIdAndTheaterIdAndShowTimeId)
}

// Movie Schedule
func (h *Handler) addMovieSchedule(ctx *gin.Context) {
	var movieSchedule MovieSchedule
	if err := ctx.BindJSON(&movieSchedule); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}

	err := h.svc.AddMovieSchedule(ctx, movieSchedule)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}

	h.response(ctx, http.StatusOK, "Movie schedule added successfully")
}

func (h *Handler) updateMovieSchedule(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	var updateData MovieSchedule
	if err := ctx.BindJSON(&updateData); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}

	err = h.svc.UpdateMovieSchedule(ctx, id, updateData)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}

	h.response(ctx, http.StatusOK, "Movie schedule updated successfully")
}

func (h *Handler) getAllMovieSchedules(ctx *gin.Context) {
	movieSchedules, err := h.svc.GetAllMovieSchedules(ctx)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}

	h.responseWithData(ctx, http.StatusOK, "movie schedule details retrieved successfully", movieSchedules)
}

func (h *Handler) getMovieScheduleByMovieID(ctx *gin.Context) {
	movieIDstr := ctx.Query("movieid")
	movieID, err := strconv.Atoi(movieIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	movieSchedules, err := h.svc.GetMovieScheduleByMovieID(ctx, movieID)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}

	h.responseWithData(ctx, http.StatusOK, "movie schedule details retrieved successfully", movieSchedules)
}

func (h *Handler) getMovieScheduleByTheaterID(ctx *gin.Context) {
	theaterIDstr := ctx.Query("theaterid")
	theaterID, err := strconv.Atoi(theaterIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	movieSchedules, err := h.svc.GetMovieScheduleByTheaterID(ctx, theaterID)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}

	h.responseWithData(ctx, http.StatusOK, "movie schedule details retrieved successfully", movieSchedules)
}

func (h *Handler) getMovieScheduleByMovieIdAndTheaterId(ctx *gin.Context) {
	movieIDstr := ctx.Query("movieid")
	theaterIDstr := ctx.Query("theaterid")
	movieID, err := strconv.Atoi(movieIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	theaterID, err := strconv.Atoi(theaterIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	movieSchedules, err := h.svc.GetMovieScheduleByMovieIdAndTheaterId(ctx, movieID, theaterID)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}

	h.responseWithData(ctx, http.StatusOK, "movie schedule details retrieved successfully", movieSchedules)
}

func (h *Handler) getMovieScheduleByMovieIdAndShowTimeId(ctx *gin.Context) {
	movieIDstr := ctx.Query("movieid")
	showTimeIDstr := ctx.Query("showtimeid")
	showTimeID, err := strconv.Atoi(showTimeIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	movieID, err := strconv.Atoi(movieIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	movieSchedules, err := h.svc.GetMovieScheduleByMovieIdAndShowTimeId(ctx, movieID, showTimeID)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "movie schedule details retrieved successfully", movieSchedules)
}

func (h *Handler) getMovieScheduleByTheaterIdAndShowTimeId(ctx *gin.Context) {
	theaterIDstr := ctx.Query("theaterid")
	showTimeIDstr := ctx.Query("showtimeid")
	showTimeID, err := strconv.Atoi(showTimeIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	theaterID, err := strconv.Atoi(theaterIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	movieSchedules, err := h.svc.GetMovieScheduleByTheaterIdAndShowTimeId(ctx, theaterID, showTimeID)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}

	h.responseWithData(ctx, http.StatusOK, "movie schedule details retrieved successfully", movieSchedules)
}

func (h *Handler) getMovieScheduleByID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	movieSchedule, err := h.svc.GetMovieScheduleByID(ctx, id)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "movie schedule details retrieved successfully", movieSchedule)
}

func (h *Handler) deleteMovieScheduleById(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	err = h.svc.DeleteMovieScheduleById(ctx, id)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}

	h.response(ctx, http.StatusOK, "movie schedule deleted successfully")
}

func (h *Handler) deleteMovieScheduleByMovieIdAndTheaterId(ctx *gin.Context) {
	movieIDstr := ctx.Query("movieid")
	theaterIDstr := ctx.Query("theaterid")
	movieID, err := strconv.Atoi(movieIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	theaterID, err := strconv.Atoi(theaterIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	err = h.svc.DeleteMovieScheduleByMovieIdAndTheaterId(ctx, movieID, theaterID)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}

	h.response(ctx, http.StatusOK, "movie schedule deleted successfully")
}

func (h *Handler) deleteMovieScheduleByMovieIdAndTheaterIdAndShowTimeId(ctx *gin.Context) {
	movieIDstr := ctx.Query("movieid")
	theaterIDstr := ctx.Query("theaterid")
	showTimeIDstr := ctx.Query("showtimeid")
	movieID, err := strconv.Atoi(movieIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	theaterID, err := strconv.Atoi(theaterIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	showTimeID, err := strconv.Atoi(showTimeIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	err = h.svc.DeleteMovieScheduleByMovieIdAndTheaterIdAndShowTimeId(ctx, movieID, theaterID, showTimeID)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}

	h.response(ctx, http.StatusOK, "movie schedule deleted successfully")
}

// Admin
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

func (h *Handler) getAdminProfile(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	admin, err := h.svc.GetAdminProfile(ctx, id)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get admin profile successfull", admin)

}

func (h *Handler) updateAdminProfile(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	admin := &AdminProfileDetails{}
	if err := ctx.ShouldBindJSON(&admin); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	err = h.svc.UpdateAdminProfile(ctx, id, *admin)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "update admin profile successfull")

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
		if formattedError == "record not found" {
			h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		} else {
			h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		}
		return
	}
	h.response(ctx, http.StatusOK, "theater deleted successfully")
}

func (h *Handler) deleteTheaterByName(ctx *gin.Context) {
	theaterName := ctx.DefaultQuery("name", "")
	err := h.svc.DeleteTheaterByName(ctx, theaterName)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		if formattedError == "record not found" {
			h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		} else {
			h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		}
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
		if formattedError == "record not found" {
			h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		} else {
			h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		}
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
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list movies succesfully", movies)
}

// Theater-types
func (h *Handler) listTheaterTypes(ctx *gin.Context) {
	theaterTypes, err := h.svc.ListTheaterTypes(ctx)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list theater-types succesfully", theaterTypes)
}

//Screen types

func (h *Handler) listScreenTypes(ctx *gin.Context) {
	screenTypes, err := h.svc.ListScreenTypes(ctx)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list screen-types successfully", screenTypes)
}

// Seat categories
func (h *Handler) listSeatCategories(ctx *gin.Context) {
	seatCategories, err := h.svc.ListSeatCategories(ctx)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list seat-categories successfully", seatCategories)
}

// Theater screen
func (h *Handler) addTheaterScreen(ctx *gin.Context) {
	theaterScreen := &TheaterScreen{}
	if err := ctx.ShouldBindJSON(&theaterScreen); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	err := h.svc.AddTheaterScreen(ctx, *theaterScreen)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "theater screen added successfully")
}

func (h *Handler) deleteTheaterScreenByID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	err = h.svc.DeleteTheaterScreenByID(ctx, id)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "theater screen deleted successfully")
}

func (h *Handler) deleteTheaterScreenByNumber(ctx *gin.Context) {
	theaterscreen := &TheaterScreen{}
	if err := ctx.ShouldBindJSON(&theaterscreen); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	err := h.svc.DeleteTheaterScreenByNumber(ctx, theaterscreen.TheaterID, theaterscreen.ScreenNumber)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "theater screen deleted successfully")
}

func (h *Handler) getTheaterScreenByID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	theaterScreen, err := h.svc.GetTheaterScreenByID(ctx, id)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get theater screen details successfully", theaterScreen)
}

func (h *Handler) getTheaterScreenByNumber(ctx *gin.Context) {
	theaterIDstr := ctx.DefaultQuery("theaterID", "")
	screenNumberstr := ctx.DefaultQuery("screenNumber", "")
	theaterID, err := strconv.Atoi(theaterIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	screenNumber, err := strconv.Atoi(screenNumberstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	theaterScreen, err := h.svc.GetTheaterScreenByNumber(ctx, theaterID, screenNumber)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get theater screen details successfully", theaterScreen)
}

func (h *Handler) updateTheaterScreen(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	theaterScreen := &TheaterScreen{}
	if err := ctx.ShouldBindJSON(&theaterScreen); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	err = h.svc.UpdateTheaterScreen(ctx, id, *theaterScreen)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "theater screen updated successfully")
}

func (h *Handler) listTheaterScreens(ctx *gin.Context) {
	theaterScreen := &TheaterScreen{}
	if err := ctx.ShouldBindJSON(&theaterScreen); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	theaterScreens, err := h.svc.ListTheaterScreens(ctx, theaterScreen.TheaterID)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list theater screens successfully", theaterScreens)
}

// Show times
func (h *Handler) addShowtime(ctx *gin.Context) {
	showtime := &Showtime{}
	if err := ctx.ShouldBindJSON(&showtime); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	err := h.svc.AddShowtime(ctx, *showtime)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "showtime added successfully")
}

func (h *Handler) deleteShowtimeByID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	err = h.svc.DeleteShowtimeByID(ctx, id)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "showtime deleted successfully")
}

func (h *Handler) deleteShowtimeByDetails(ctx *gin.Context) {
	movieIDStr := ctx.DefaultQuery("movie_id", "")
	screenIDStr := ctx.DefaultQuery("screen_id", "")
	showDateStr := ctx.DefaultQuery("show_date", "")
	showTimeStr := ctx.DefaultQuery("show_time", "")

	movieID, _ := strconv.Atoi(movieIDStr)
	screenID, _ := strconv.Atoi(screenIDStr)
	showDate, _ := time.Parse(time.RFC3339, showDateStr)
	showTime, _ := time.Parse(time.RFC3339, showTimeStr)

	err := h.svc.DeleteShowtimeByDetails(ctx, movieID, screenID, showDate, showTime)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "showtime deleted successfully")
}

func (h *Handler) getShowtimeByID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	showtime, err := h.svc.GetShowtimeByID(ctx, id)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get showtime details successfully", showtime)
}

func (h *Handler) getShowtimeByDetails(ctx *gin.Context) {
	movieIDStr := ctx.DefaultQuery("movie_id", "")
	screenIDStr := ctx.DefaultQuery("screen_id", "")
	showDateStr := ctx.DefaultQuery("show_date", "")
	showTimeStr := ctx.DefaultQuery("show_time", "")

	movieID, _ := strconv.Atoi(movieIDStr)
	screenID, _ := strconv.Atoi(screenIDStr)
	showDate, _ := time.Parse(time.RFC3339, showDateStr)
	showTime, _ := time.Parse(time.RFC3339, showTimeStr)

	showtime, err := h.svc.GetShowtimeByDetails(ctx, movieID, screenID, showDate, showTime)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get showtime details successfully", showtime)
}

func (h *Handler) updateShowtime(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	showtime := &Showtime{}
	if err := ctx.ShouldBindJSON(&showtime); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	err = h.svc.UpdateShowtime(ctx, id, *showtime)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.response(ctx, http.StatusOK, "showtime updated successfully")
}

func (h *Handler) listShowtimes(ctx *gin.Context) {
	movieIDStr := ctx.DefaultQuery("movie_id", "")
	movieID, _ := strconv.Atoi(movieIDStr)

	showtimes, err := h.svc.ListShowtimes(ctx, movieID)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list showtimes successfully", showtimes)
}
