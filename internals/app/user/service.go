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
	GetMovieByNameAndLanguage(ctx context.Context, name, language string) (*Movie, error)
	// Theater
	ListAllTheaters(ctx context.Context) (interface{}, error)
	GetTheaterByID(ctx context.Context, id int) (*TheaterWithTypeResponse, error)
	GetTheatersByCity(ctx context.Context, city string) ([]TheaterWithTypeResponse, error)
	GetTheatersByName(ctx context.Context, name string) ([]TheaterWithTypeResponse, error)
	GetTheatersAndMovieScheduleByMovieName(ctx context.Context, movieName string) ([]TheatersAndMovieScheduleResponse, error)
	GetScreensAndMovieSchedulesByTheaterID(ctx context.Context, id int) (*TheaterResponse, error)
	ListShowTimeByTheaterID(ctx context.Context, id int) (*ListShowTimeResponse, error)
	ListShowTimeByTheaterIDandMovieID(ctx context.Context, theaterId int, movieId int) (*ListShowTimeByTheaterAndMovie, error)
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

func (s *service) GetMovieByNameAndLanguage(ctx context.Context, name, language string) (*Movie, error) {
	res, err := s.movieBooking.GetMovieByNameAndLanguage(ctx, &movie_booking.GetMovieByNameAndLanguageRequest{
		Name:     name,
		Language: language,
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

// Theater
func (s *service) ListAllTheaters(ctx context.Context) (interface{}, error) {
	response, err := s.theaterClient.ListTheaters(ctx, &movie_booking.ListTheatersRequest{})
	if err != nil {
		return nil, err
	}
	theaterResponses := []TheaterWithTypeResponse{}

	for _, theater := range response.Theaters {
		theaterResponse := TheaterWithTypeResponse{
			ID:              int(theater.TheaterId),
			Name:            theater.Name,
			Place:           theater.Place,
			City:            theater.City,
			District:        theater.District,
			State:           theater.State,
			OwnerID:         int(theater.OwnerId),
			NumberOfScreens: int(theater.NumberOfScreens),
			TheaterType: TheaterTypeResponse{
				ID:              int(theater.TheaterType.Id),
				TheaterTypeName: theater.TheaterType.TheaterTypeName,
			},
		}
		theaterResponses = append(theaterResponses, theaterResponse)
	}

	return theaterResponses, nil
}

func (s *service) GetTheaterByID(ctx context.Context, id int) (*TheaterWithTypeResponse, error) {
	response, err := s.theaterClient.GetTheaterByID(ctx, &movie_booking.GetTheaterByIDRequest{
		TheaterId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	theaterResponses := TheaterWithTypeResponse{
		ID:              int(response.Theater.TheaterId),
		Name:            response.Theater.Name,
		Place:           response.Theater.Place,
		City:            response.Theater.City,
		District:        response.Theater.District,
		State:           response.Theater.State,
		OwnerID:         int(response.Theater.OwnerId),
		NumberOfScreens: int(response.Theater.NumberOfScreens),
		TheaterType: TheaterTypeResponse{
			ID:              int(response.Theater.TheaterType.Id),
			TheaterTypeName: response.Theater.TheaterType.TheaterTypeName,
		},
	}
	return &theaterResponses, nil
}

func (s *service) GetTheatersByCity(ctx context.Context, city string) ([]TheaterWithTypeResponse, error) {
	response, err := s.theaterClient.GetTheatersByCity(ctx, &movie_booking.GetTheatersByCityRequest{
		City: city,
	})
	if err != nil {
		return nil, err
	}
	theaterResponses := []TheaterWithTypeResponse{}

	for _, theater := range response.Theater {
		theaterResponse := TheaterWithTypeResponse{
			ID:              int(theater.TheaterId),
			Name:            theater.Name,
			Place:           theater.Place,
			City:            theater.City,
			District:        theater.District,
			State:           theater.State,
			OwnerID:         int(theater.OwnerId),
			NumberOfScreens: int(theater.NumberOfScreens),
			TheaterType: TheaterTypeResponse{
				ID:              int(theater.TheaterType.Id),
				TheaterTypeName: theater.TheaterType.TheaterTypeName,
			},
		}
		theaterResponses = append(theaterResponses, theaterResponse)
	}

	return theaterResponses, nil
}

func (s *service) GetTheatersAndMovieScheduleByMovieName(ctx context.Context, movieName string) ([]TheatersAndMovieScheduleResponse, error) {
	response, err := s.theaterClient.GetTheatersAndMovieScheduleByMovieName(ctx, &movie_booking.GetTheatersAndMovieScheduleByMovieNameRequest{
		Name: movieName,
	})
	if err != nil {
		return nil, err
	}
	theaterResponses := []TheatersAndMovieScheduleResponse{}

	for _, res := range response.MovieScedule {
		theaterResponse := TheatersAndMovieScheduleResponse{
			ID:         int(res.Id),
			MovieID:    int(res.MovieId),
			TheaterID:  int(res.TheaterId),
			ShowtimeID: int(res.ShowtimeId),
			Movie: Movie{
				Title:       res.Movie.Title,
				Description: res.Movie.Description,
				Duration:    int(res.Movie.Duration),
				Genre:       res.Movie.Genre,
				ReleaseDate: res.Movie.ReleaseDate,
				Rating:      float64(res.Movie.Rating),
				Language:    res.Movie.Language,
			},
			Theater: Theater{
				ID:              uint(res.TheaterId),
				Name:            res.Theater.Name,
				Place:           res.Theater.Place,
				City:            res.Theater.City,
				District:        res.Theater.District,
				State:           res.Theater.State,
				OwnerID:         uint(res.Theater.OwnerId),
				NumberOfScreens: int(res.Theater.NumberOfScreens),
				TheaterTypeID:   int(res.Theater.TheaterTypeId),
			},
			Showtime: Showtime{
				ID:       uint(res.ShowTime.Id),
				MovieID:  int(res.ShowTime.MovieId),
				ScreenID: int(res.ShowTime.ScreenId),
				ShowDate: res.ShowTime.ShowDate.AsTime(),
				ShowTime: res.ShowTime.ShowTime.AsTime(),
			},
		}
		theaterResponses = append(theaterResponses, theaterResponse)
	}
	return theaterResponses, nil
}

func (s *service) GetTheatersByName(ctx context.Context, name string) ([]TheaterWithTypeResponse, error) {
	response, err := s.theaterClient.GetTheaterByName(ctx, &movie_booking.GetTheaterByNameRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	theaterResponses := []TheaterWithTypeResponse{}

	for _, theater := range response.Theater {
		theaterResponse := TheaterWithTypeResponse{
			ID:              int(theater.TheaterId),
			Name:            theater.Name,
			Place:           theater.Place,
			City:            theater.City,
			District:        theater.District,
			State:           theater.State,
			OwnerID:         int(theater.OwnerId),
			NumberOfScreens: int(theater.NumberOfScreens),
			TheaterType: TheaterTypeResponse{
				ID:              int(theater.TheaterType.Id),
				TheaterTypeName: theater.TheaterType.TheaterTypeName,
			},
		}
		theaterResponses = append(theaterResponses, theaterResponse)
	}
	return theaterResponses, nil
}
func (s *service) GetScreensAndMovieSchedulesByTheaterID(ctx context.Context, id int) (*TheaterResponse, error) {
	response, err := s.theaterClient.GetScreensAndMovieScedulesByTheaterID(ctx, &movie_booking.GetScreensAndMovieScedulesByTheaterIdRequest{
		TheaterId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	theater := Theater{
		ID:              uint(id),
		Name:            response.Theater.Name,
		Place:           response.Theater.Place,
		City:            response.Theater.City,
		District:        response.Theater.District,
		State:           response.Theater.State,
		NumberOfScreens: int(response.Theater.NumberOfScreens),
	}
	theaterType := TheaterType{
		ID:              int(response.TheaterType.Id),
		TheaterTypeName: response.TheaterType.TheaterTypeName,
	}
	var movieSchedules []MovieSchedule
	for _, resSchedule := range response.MovieSchedule {
		movieSchedule := MovieSchedule{
			ID:         int(resSchedule.Id),
			MovieID:    int(resSchedule.MovieId),
			TheaterID:  int(resSchedule.TheaterId),
			ShowtimeID: int(resSchedule.ShowtimeId),
			Movie: Movie{
				Title:       resSchedule.Movie.Title,
				Description: resSchedule.Movie.Description,
				Duration:    int(resSchedule.Movie.Duration),
				Genre:       resSchedule.Movie.Genre,
				ReleaseDate: resSchedule.Movie.ReleaseDate,
				Rating:      float64(resSchedule.Movie.Rating),
				Language:    resSchedule.Movie.Language,
			},
			Showtime: Showtime{
				ID:       uint(resSchedule.ShowTime.Id),
				MovieID:  int(resSchedule.ShowTime.MovieId),
				ScreenID: int(resSchedule.ShowTime.ScreenId),
				ShowDate: resSchedule.ShowTime.ShowDate.AsTime(),
				ShowTime: resSchedule.ShowTime.ShowTime.AsTime(),
			},
		}
		movieSchedules = append(movieSchedules, movieSchedule)
	}

	var theaterScreens []TheaterScreen
	for _, resScreen := range response.TheaterScreen {
		theaterScreen := TheaterScreen{
			ID:           uint(resScreen.ID),
			TheaterID:    int(resScreen.TheaterID),
			ScreenNumber: int(resScreen.ScreenNumber),
			SeatCapacity: int(resScreen.SeatCapacity),
			ScreenTypeID: int(resScreen.ScreenTypeID),
			ScreenType: ScreenType{
				ID:             int(resScreen.ScreenType.Id),
				ScreenTypeName: resScreen.ScreenType.ScreenTypeName,
			},
		}
		theaterScreens = append(theaterScreens, theaterScreen)
	}
	return &TheaterResponse{
		ID:              int(theater.ID),
		Name:            theater.Name,
		Place:           theater.Place,
		City:            theater.City,
		District:        theater.District,
		State:           theater.State,
		NumberOfScreens: theater.NumberOfScreens,
		TheaterType:     theaterType,
		MovieSchedules:  movieSchedules,
		TheaterScreens:  theaterScreens,
	}, nil
}

func (s *service) ListShowTimeByTheaterID(ctx context.Context, id int) (*ListShowTimeResponse, error) {
	response, err := s.theaterClient.ListShowTimeByTheaterID(ctx, &movie_booking.ListShowTimeByTheaterIdRequest{
		TheaterId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	theater := Theater{
		ID:              uint(response.Theater.TheaterId),
		Name:            response.Theater.Name,
		Place:           response.Theater.Place,
		City:            response.Theater.City,
		District:        response.Theater.District,
		State:           response.Theater.State,
		OwnerID:         uint(response.Theater.OwnerId),
		NumberOfScreens: int(response.Theater.NumberOfScreens),
		TheaterTypeID:   int(response.Theater.TheaterTypeId),
	}
	showTimeResponse := []ShowtimeResponse{}

	for _, res := range response.ShowTime {
		showtime := ShowtimeResponse{
			ID:       uint(res.Id),
			MovieID:  int(res.MovieId),
			ScreenID: int(res.ScreenId),
			ShowDate: res.ShowDate.AsTime(),
			ShowTime: res.ShowTime.AsTime(),
			Movie: Movie{
				Title:       res.Movie.Title,
				Description: res.Movie.Description,
				Duration:    int(res.Movie.Duration),
				Genre:       res.Movie.Genre,
				ReleaseDate: res.Movie.ReleaseDate,
				Rating:      float64(res.Movie.Rating),
				Language:    res.Movie.Language,
			},
			TheaterScreen: TheaterScreen{
				ID:           uint(res.TheaterScreen.ID),
				TheaterID:    int(res.TheaterScreen.TheaterID),
				ScreenNumber: int(res.TheaterScreen.ScreenNumber),
				SeatCapacity: int(res.TheaterScreen.SeatCapacity),
				ScreenTypeID: int(res.TheaterScreen.ScreenTypeID),
			},
		}
		showTimeResponse = append(showTimeResponse, showtime)
	}
	return &ListShowTimeResponse{
		Theater:          theater,
		ShowtimeResponse: showTimeResponse,
	}, nil
}

func (s *service) ListShowTimeByTheaterIDandMovieID(ctx context.Context, theaterId int, movieId int) (*ListShowTimeByTheaterAndMovie, error) {
	response, err := s.theaterClient.ListShowTimeByTheaterIDandMovieID(ctx, &movie_booking.ListShowTimeByTheaterIdandMovieIdRequest{
		TheaterId: int32(theaterId),
		MovieId:   int32(movieId),
	})
	if err != nil {
		return nil, err
	}
	movie := Movie{
		Title:       response.Movie.Title,
		Description: response.Movie.Description,
		Duration:    int(response.Movie.Duration),
		Genre:       response.Movie.Genre,
		ReleaseDate: response.Movie.ReleaseDate,
		Rating:      float64(response.Movie.Rating),
		Language:    response.Movie.Language,
	}
	theater := Theater{
		ID:              uint(response.Theater.TheaterId),
		Name:            response.Theater.Name,
		Place:           response.Theater.Place,
		City:            response.Theater.City,
		District:        response.Theater.District,
		State:           response.Theater.State,
		OwnerID:         uint(response.Theater.OwnerId),
		NumberOfScreens: int(response.Theater.NumberOfScreens),
		TheaterTypeID:   int(response.Theater.TheaterTypeId),
	}
	showTimeResponse := []ShowtimeResponse{}

	for _, res := range response.ShowTime {
		showtime := ShowtimeResponse{
			ID:       uint(res.Id),
			MovieID:  int(res.MovieId),
			ScreenID: int(res.ScreenId),
			ShowDate: res.ShowDate.AsTime(),
			ShowTime: res.ShowTime.AsTime(),
			TheaterScreen: TheaterScreen{
				ID:           uint(res.TheaterScreen.ID),
				TheaterID:    int(res.TheaterScreen.TheaterID),
				ScreenNumber: int(res.TheaterScreen.ScreenNumber),
				SeatCapacity: int(res.TheaterScreen.SeatCapacity),
				ScreenTypeID: int(res.TheaterScreen.ScreenTypeID),
			},
		}
		showTimeResponse = append(showTimeResponse, showtime)
	}
	return &ListShowTimeByTheaterAndMovie{
		Movie:            movie,
		Theater:          theater,
		ShowtimeResponse: showTimeResponse,
	}, nil
}
