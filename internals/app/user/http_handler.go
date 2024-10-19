package user

import (
	"errors"
	"fmt"
	"io"
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
	r.POST("/register/validate", h.registerValidate)
	r.POST("/login", h.logIn)
	r.POST("/forgot/password", h.forgotPassword)
	r.POST("/reset/password", h.resetPassword)

	r.GET("/movies", h.listAllMovies)
	r.GET("/movie/:id", h.getMovieDetailsByID)
	r.GET("/movie/name", h.getMovieByName)
	r.GET("/movie/genre", h.getMoviesByGenre)
	r.GET("/movie/language", h.getMoviesByLanguage)
	r.GET("/movie/name/language", h.getMovieByNameAndLanguage)

	r.GET("/theaters", h.listAllTheaters)
	r.GET("/theater/:id", h.getTheaterByID)
	r.GET("/theaters/name", h.getTheatersByName)
	r.GET("/theaters/city", h.getTheatersByCity)
	r.GET("/theaters/movie/name", h.getTheatersAndMovieScheduleByMovieName)
	r.GET("/theater/details/:id", h.getScreensAndMovieScedulesByTheaterID)
	r.GET("/theater/showtime/:id", h.listShowTimeByTheaterID)
	r.GET("/theaters/:theater_id/movies/:movie_id/showtimes", h.listShowTimeByTheaterIDandMovieID)
	r.GET("/theater/movie/showdate/showtimes/:movie_id", h.listShowtimeByMovieIdAndShowDate)

	r.GET("/theater/screen/seats/:screen_id", h.listSeatsbyScreenID)
	r.GET("/theater/screen/seat/:seat_id", h.getSeatBySeatID)

	auth := r.Use(h.authHandler.UserAuthMiddleware())
	auth.GET("/profile", h.getProfile)
	auth.PUT("/profile/:id", h.updateUserProfile)
	//Booking
	auth.POST("/booking", h.createBooking)
	auth.GET("/booking/:id", h.getBookingByID)
	auth.GET("/booking/user/:user_id", h.listBookingsByUser)
	// Payment
	auth.GET("/payment/status/:transaction_id", h.getTransactionStatus)
	auth.POST("/payment/:booking_id", h.processPayment)
	auth.POST("/webhook/razorpay", h.handleRazorpayWebhook)
}

// Payment
func (h *Handler) handleRazorpayWebhook(ctx *gin.Context) {
	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.responseWithError(ctx, http.StatusBadRequest, errors.New("failed to read request body"))
		return
	}

	if err := h.svc.HandleRazorpayWebhook(ctx, payload); err != nil {
		h.responseWithError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to handle webhook: %v", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

func (h *Handler) processPayment(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("Authorization")
	userId, err := h.svc.GetUserIDFromToken(ctx, authorization)
	if err != nil {
		h.responseWithError(ctx, http.StatusUnauthorized, errors.New("unauthorized: Invalid token or user ID extraction failed"))
		return
	}
	bookingIdstr := ctx.Param("booking_id")
	bookingId, err := strconv.Atoi(bookingIdstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	transaction, err := h.svc.ProcessPayment(ctx, bookingId, userId)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "payment processed succesfully",
		"transaction": transaction,
		"redirect":    fmt.Sprintf("http://localhost:8080/gateway/user/webhook/razorpay/%d", transaction.TransactionID),
	})
}

func (h *Handler) getTransactionStatus(ctx *gin.Context) {
	idstr := ctx.Param("transaction_id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}

	status, err := h.svc.GetTransactionStatus(ctx, id)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}

	h.responseWithData(ctx, http.StatusOK, "transaction status retrieved", status)
}

// Booking
func (h *Handler) listBookingsByUser(ctx *gin.Context) {
	idstr := ctx.Param("user_id")
	userId, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	bookings, err := h.svc.ListBookingsByUser(ctx, userId)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get booking details succesfull", bookings)

}

func (h *Handler) getBookingByID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	bookings, err := h.svc.GetBookingByID(ctx, id)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get booking details succesfull", bookings)

}
func (h *Handler) createBooking(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("Authorization")
	userId, err := h.svc.GetUserIDFromToken(ctx, authorization)
	if err != nil {
		h.responseWithError(ctx, http.StatusUnauthorized, errors.New("unauthorized: Invalid token or user ID extraction failed"))
		return
	}
	bookingReq := &CreateBookingRequest{}
	if err := ctx.ShouldBindJSON(bookingReq); err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	bookingReq.UserID = userId
	booking, err := h.svc.CreateBooking(ctx, *bookingReq)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusBadRequest, errors.New(formattedError))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "booking successful",
		"booking":  booking,
		"redirect": fmt.Sprintf("http://localhost:8080/gateway/user/payment/%d", booking.BookingID),
	})
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

// Movies
func (h *Handler) listAllMovies(ctx *gin.Context) {
	movies, err := h.svc.ListAllMovies(ctx)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list all movies successfully", movies)
}

func (h *Handler) getMovieDetailsByID(ctx *gin.Context) {
	movieIDstr := ctx.Param("id")
	movieID, err := strconv.Atoi(movieIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	movie, err := h.svc.GetMovieDetailsByID(ctx, movieID)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get movie details successfully", movie)
}

func (h *Handler) getMovieByName(ctx *gin.Context) {
	movieName := ctx.Query("name")
	movie, err := h.svc.GetMovieByName(ctx, movieName)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get movie by name successfully", movie)
}

func (h *Handler) getMoviesByGenre(ctx *gin.Context) {
	genre := ctx.Query("genre")
	movies, err := h.svc.GetMoviesByGenre(ctx, genre)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get movies by genre successfully", movies)
}

func (h *Handler) getMoviesByLanguage(ctx *gin.Context) {
	language := ctx.Query("language")
	movies, err := h.svc.GetMoviesByLanguage(ctx, language)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get movies by language successfully", movies)
}

func (h *Handler) getMovieByNameAndLanguage(ctx *gin.Context) {
	name := ctx.Query("name")
	language := ctx.Query("language")

	movie, err := h.svc.GetMovieByNameAndLanguage(ctx, name, language)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}

	h.responseWithData(ctx, http.StatusOK, "get movie by name and language successfully", movie)
}

// Theaters
func (h *Handler) listAllTheaters(ctx *gin.Context) {
	theaters, err := h.svc.ListAllTheaters(ctx)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list all theaters successfully", theaters)
}

func (h *Handler) getTheaterByID(ctx *gin.Context) {
	theaterIDstr := ctx.Param("id")
	theaterID, err := strconv.Atoi(theaterIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	theater, err := h.svc.GetTheaterByID(ctx, theaterID)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get theater by ID successfully", theater)
}

func (h *Handler) getTheatersByCity(ctx *gin.Context) {
	city := ctx.Query("city")
	theaters, err := h.svc.GetTheatersByCity(ctx, city)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get theaters by city successfully", theaters)
}

func (h *Handler) getTheatersByName(ctx *gin.Context) {
	theaterName := ctx.Query("name")
	theaters, err := h.svc.GetTheatersByName(ctx, theaterName)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get theaters by name successfully", theaters)
}

func (h *Handler) getTheatersAndMovieScheduleByMovieName(ctx *gin.Context) {
	movieName := ctx.Query("movie_name")
	theaters, err := h.svc.GetTheatersAndMovieScheduleByMovieName(ctx, movieName)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get theaters by movie name successfully", theaters)
}

func (h *Handler) getScreensAndMovieScedulesByTheaterID(ctx *gin.Context) {
	theaterIDstr := ctx.Param("id")
	theaterID, err := strconv.Atoi(theaterIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	screens, err := h.svc.GetScreensAndMovieSchedulesByTheaterID(ctx, theaterID)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "get screens and movie schedules by theater ID successfully", screens)
}

func (h *Handler) listShowTimeByTheaterID(ctx *gin.Context) {
	theaterIDstr := ctx.Param("id")
	theaterID, err := strconv.Atoi(theaterIDstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	showtimes, err := h.svc.ListShowTimeByTheaterID(ctx, theaterID)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list showtimes by theater ID successfully", showtimes)
}

func (h *Handler) listShowTimeByTheaterIDandMovieID(ctx *gin.Context) {
	theaterIDstr := ctx.Param("theater_id")
	movieIDstr := ctx.Param("movie_id")
	theaterID, err := strconv.Atoi(theaterIDstr)
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
	showtimes, err := h.svc.ListShowTimeByTheaterIDandMovieID(ctx, theaterID, movieID)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list showtimes by theater ID and movie ID successfully", showtimes)
}

func (h *Handler) listSeatsbyScreenID(ctx *gin.Context) {
	idstr := ctx.Param("screen_id")
	screenId, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	seats, err := h.svc.ListSeatsbyScreenID(ctx, screenId)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list seats by screen id successfully", seats)

}

func (h *Handler) getSeatBySeatID(ctx *gin.Context) {
	idstr := ctx.Param("seat_id")
	seatId, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	seat, err := h.svc.GetSeatBySeatID(ctx, seatId)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list seats by seat id successfully", seat)
}

func (h *Handler) listShowtimeByMovieIdAndShowDate(ctx *gin.Context) {
	idstr := ctx.Param("movie_id")
	movieId, err := strconv.Atoi(idstr)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusInternalServerError, errors.New(formattedError))
		return
	}
	showDateStr := ctx.DefaultQuery("show_date", "")
	showDate, _ := time.Parse(time.RFC3339, showDateStr)
	showtimes, err := h.svc.ListShowtimeByMovieIdAndShowDate(ctx, showDate, movieId)
	if err != nil {
		formattedError := ExtractErrorMessage(err)
		h.responseWithError(ctx, http.StatusNotFound, errors.New(formattedError))
		return
	}
	h.responseWithData(ctx, http.StatusOK, "list showtimes by movie_id and show date successfully", showtimes)
}
