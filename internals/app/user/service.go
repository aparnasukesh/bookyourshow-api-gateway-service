package user

import (
	"context"
	"strings"

	"github.com/aparnasukesh/inter-communication/auth"
	"github.com/aparnasukesh/inter-communication/movie_booking"
	"github.com/aparnasukesh/inter-communication/user_admin"
)

type Service interface {
	Register(ctx context.Context, signUpData *User) error
	RegisterValidate(ctx context.Context, userData *User) error
	Login(ctx context.Context, loginData *User) (string, error)
	GetUserIDFromToken(ctx context.Context, authorization string) (int, error)
	GetProfile(ctx context.Context, userId int) (*UserProfileDetails, error)
	UpdateUserProfile(ctx context.Context, id int, user UserProfileDetails) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, data ResetPassword) error
	// Movies
	ListAllMovies(ctx context.Context) ([]Movie, error)
	GetMovieDetailsByID(ctx context.Context, id int) (*Movie, error)
	GetMovieByName(ctx context.Context, name string) (*Movie, error)
	GetMoviesByGenre(ctx context.Context, genre string) ([]Movie, error)
	GetMoviesByLanguage(ctx context.Context, language string) ([]Movie, error)
	// Theater
	ListAllTheaters(ctx context.Context) ([]Theater, error)
	GetTheaterByID(ctx context.Context, id int) (*Theater, error)
	GetTheatersByCity(ctx context.Context, city string) ([]Theater, error)
	GetTheatersByName(ctx context.Context, name string) ([]Theater, error)
	GetTheatersByMovieName(ctx context.Context, movieName string) ([]Theater, error)
	GetScreensAndMovieSchedulesByTheaterID(ctx context.Context, id int) (*Theater, error)
	ListShowTimeByTheaterID(ctx context.Context, id int) ([]Showtime, error)
	ListShowTimeByTheaterIDandMovieID(ctx context.Context, theaterId, movieId int) ([]Showtime, error)
}

type service struct {
	userAdmin     user_admin.UserServiceClient
	auth          auth.JWT_TokenServiceClient
	movieBooking  movie_booking.MovieServiceClient
	theaterClient movie_booking.TheatreServiceClient
}

func NewService(pb user_admin.UserServiceClient, auth auth.JWT_TokenServiceClient, movieBooking movie_booking.MovieServiceClient, theaterClient movie_booking.TheatreServiceClient) Service {
	return &service{
		userAdmin:     pb,
		auth:          auth,
		movieBooking:  movieBooking,
		theaterClient: theaterClient,
	}
}

func (s *service) Register(ctx context.Context, signUpData *User) error {
	reqData := user_admin.RegisterUserRequest{
		Username:  signUpData.Username,
		Password:  signUpData.Password,
		Phone:     signUpData.PhoneNumber,
		Email:     signUpData.Email,
		FirstName: signUpData.FirstName,
		LastName:  signUpData.LastName,
		Gender:    signUpData.Gender,
	}
	if _, err := s.userAdmin.RegisterUser(ctx, &reqData); err != nil {
		return err
	}
	return nil
}

func (s *service) RegisterValidate(ctx context.Context, userData *User) error {
	req := user_admin.ValidateUserRequest{
		Email: userData.Email,
		Otp:   userData.Otp,
	}
	if _, err := s.userAdmin.ValidateUser(ctx, &req); err != nil {
		return err
	}
	return nil
}
func (s *service) Login(ctx context.Context, loginData *User) (string, error) {
	req := user_admin.LoginUserRequest{
		Email:    loginData.Email,
		Password: loginData.Password,
	}
	res, err := s.userAdmin.LoginUser(ctx, &req)
	if err != nil {
		return "", err
	}
	return res.Token, nil
}

func (s *service) GetUserIDFromToken(ctx context.Context, authorization string) (int, error) {
	tokenParts := strings.Split(authorization, "Bearer ")
	token := tokenParts[1]
	var userId int

	response, err := s.auth.GetUserID(ctx, &auth.GetUserIDRequest{
		Token: token,
	})
	userId = int(response.UserId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (s *service) GetProfile(ctx context.Context, userId int) (*UserProfileDetails, error) {
	response, err := s.userAdmin.GetUserProfile(ctx, &user_admin.GetProfileRequest{
		UserId: int32(userId),
	})
	if err != nil {
		return nil, err
	}
	details, err := BuildGetUserProfile(response)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func (s *service) UpdateUserProfile(ctx context.Context, id int, user UserProfileDetails) error {
	_, err := s.userAdmin.UpdateUserProfile(ctx, &user_admin.UpdateUserProfileRequest{
		UserId:      int32(id),
		Username:    user.Username,
		Phone:       user.PhoneNumber,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Gender:      user.Gender,
		DateOfBirth: user.DateOfBirth,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ForgotPassword(ctx context.Context, email string) error {
	_, err := s.userAdmin.ForgotUserPassword(ctx, &user_admin.ForgotPasswordRequest{
		Email: email,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ResetPassword(ctx context.Context, data ResetPassword) error {
	_, err := s.userAdmin.ResetUserPassword(ctx, &user_admin.ResetPasswordRequest{
		Email:       data.Email,
		Otp:         data.Otp,
		NewPassword: data.NewPassword,
	})
	if err != nil {
		return err
	}
	return nil
}

// Movies

func (s *service) ListAllMovies(ctx context.Context) ([]Movie, error) {
	response, err := s.movieBooking.ListMovies(ctx, &movie_booking.ListMoviesRequest{})
	if err != nil {
		return nil, err
	}
	movies := []Movie{}
	for _, res := range response.Movies {
		movie := Movie{
			Title:       res.Title,
			Description: res.Description,
			Duration:    int(res.Duration),
			Genre:       res.Genre,
			ReleaseDate: res.ReleaseDate,
			Rating:      float64(res.Rating),
			Language:    res.Language,
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (s *service) GetMovieByName(ctx context.Context, name string) (*Movie, error) {
	res, err := s.movieBooking.GetMovieByName(ctx, &movie_booking.GetMovieByNameRequest{
		MovieName: name,
	})
	if err != nil {
		return nil, err
	}
	movie := &Movie{
		Title:       res.Movie.Title,
		Description: res.Movie.Description,
		Duration:    int(res.Movie.Duration),
		Genre:       res.Movie.Genre,
		ReleaseDate: res.Movie.ReleaseDate,
		Rating:      float64(res.Movie.Rating),
		Language:    res.Movie.Language,
	}
	return movie, nil
}

func (s *service) GetMovieDetailsByID(ctx context.Context, id int) (*Movie, error) {
	res, err := s.movieBooking.GetMovieDetailsByID(ctx, &movie_booking.GetMovieDetailsRequest{
		MovieId: uint32(id),
	})
	if err != nil {
		return nil, err
	}
	movie := &Movie{
		Title:       res.Movie.Title,
		Description: res.Movie.Description,
		Duration:    int(res.Movie.Duration),
		Genre:       res.Movie.Genre,
		ReleaseDate: res.Movie.ReleaseDate,
		Rating:      float64(res.Movie.Rating),
		Language:    res.Movie.Language,
	}
	return movie, nil
}

func (s *service) GetMoviesByGenre(ctx context.Context, genre string) ([]Movie, error) {
	response, err := s.movieBooking.GetMoviesByGenre(ctx, &movie_booking.GetMoviesByGenreRequest{
		Genre: genre,
	})
	if err != nil {
		return nil, err
	}
	movies := []Movie{}
	for _, res := range response.Movie {
		movie := Movie{
			Title:       res.Title,
			Description: res.Description,
			Duration:    int(res.Duration),
			Genre:       res.Genre,
			ReleaseDate: res.ReleaseDate,
			Rating:      float64(res.Rating),
			Language:    res.Language,
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (s *service) GetMoviesByLanguage(ctx context.Context, language string) ([]Movie, error) {
	response, err := s.movieBooking.GetMoviesByLanguage(ctx, &movie_booking.GetMoviesByLanguageRequest{
		Language: language,
	})
	if err != nil {
		return nil, err
	}
	movies := []Movie{}
	for _, res := range response.Movie {
		movie := Movie{
			Title:       res.Title,
			Description: res.Description,
			Duration:    int(res.Duration),
			Genre:       res.Genre,
			ReleaseDate: res.ReleaseDate,
			Rating:      float64(res.Rating),
			Language:    res.Language,
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

// Theater

func (s *service) ListAllTheaters(ctx context.Context) ([]Theater, error) {
	panic("unimplemented")
}

func (s *service) GetScreensAndMovieSchedulesByTheaterID(ctx context.Context, id int) (*Theater, error) {
	panic("unimplemented")
}

func (s *service) GetTheaterByID(ctx context.Context, id int) (*Theater, error) {
	panic("unimplemented")
}

func (s *service) GetTheatersByCity(ctx context.Context, city string) ([]Theater, error) {
	panic("unimplemented")
}

func (s *service) GetTheatersByMovieName(ctx context.Context, movieName string) ([]Theater, error) {
	panic("unimplemented")
}

func (s *service) GetTheatersByName(ctx context.Context, name string) ([]Theater, error) {
	panic("unimplemented")
}

func (s *service) ListShowTimeByTheaterID(ctx context.Context, id int) ([]Showtime, error) {
	panic("unimplemented")
}

func (s *service) ListShowTimeByTheaterIDandMovieID(ctx context.Context, theaterId int, movieId int) ([]Showtime, error) {
	panic("unimplemented")
}
