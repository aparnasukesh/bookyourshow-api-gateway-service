package admin

import (
	"context"
	"strings"
	"time"

	"github.com/aparnasukesh/inter-communication/auth"
	"github.com/aparnasukesh/inter-communication/user_admin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service interface {
	Register(ctx context.Context, signUpData *Admin) error
	Login(ctx context.Context, loginData *Admin) (string, error)
	GetUserIDFromToken(ctx context.Context, authorization string) (int, error)
	GetAdminProfile(ctx context.Context, id int) (*Admin, error)
	UpdateAdminProfile(ctx context.Context, id int, admin AdminProfileDetails) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, data ResetPassword) error
	//Theater
	AddTheater(ctx context.Context, theater Theater) error
	DeleteTheaterByID(ctx context.Context, id int) error
	DeleteTheaterByName(ctx context.Context, name string) error
	GetTheaterByID(ctx context.Context, id int) (*Theater, error)
	GetTheaterByName(ctx context.Context, name string) ([]Theater, error)
	UpdateTheater(ctx context.Context, id int, theater Theater) error
	ListTheaters(ctx context.Context) ([]Theater, error)
	//Movies
	ListMovies(ctx context.Context) ([]Movie, error)
	//Theater types
	ListTheaterTypes(ctx context.Context) ([]TheaterType, error)
	//Screen type
	ListScreenTypes(ctx context.Context) ([]ScreenType, error)
	//Seat categories
	ListSeatCategories(ctx context.Context) ([]SeatCategory, error)
	//Theater screen
	AddTheaterScreen(ctx context.Context, theaterScreen TheaterScreen) error
	DeleteTheaterScreenByID(ctx context.Context, id int) error
	DeleteTheaterScreenByNumber(ctx context.Context, theaterID int, screenNumber int) error
	GetTheaterScreenByID(ctx context.Context, id int) (*TheaterScreen, error)
	GetTheaterScreenByNumber(ctx context.Context, theaterID int, screenNumber int) (*TheaterScreen, error)
	UpdateTheaterScreen(ctx context.Context, id int, theaterScreen TheaterScreen) error
	ListTheaterScreens(ctx context.Context, theaterId int) ([]TheaterScreen, error)
	//Show Time
	AddShowtime(ctx context.Context, showtime Showtime) error
	DeleteShowtimeByID(ctx context.Context, id int) error
	DeleteShowtimeByDetails(ctx context.Context, movieID int, screenID int, showDate time.Time, showTime time.Time) error
	GetShowtimeByID(ctx context.Context, id int) (*Showtime, error)
	GetShowtimeByDetails(ctx context.Context, movieID int, screenID int, showDate time.Time, showTime time.Time) (*Showtime, error)
	UpdateShowtime(ctx context.Context, id int, showtime Showtime) error
	ListShowtimes(ctx context.Context, movieID int) ([]Showtime, error)
	//Show time
	AddMovieSchedule(ctx context.Context, movieSchedule MovieSchedule) error
	UpdateMovieSchedule(ctx context.Context, id int, updateData MovieSchedule) error
	GetAllMovieSchedules(ctx context.Context) ([]MovieSchedule, error)
	GetMovieScheduleByMovieID(ctx context.Context, id int) ([]MovieSchedule, error)
	GetMovieScheduleByTheaterID(ctx context.Context, id int) ([]MovieSchedule, error)
	GetMovieScheduleByMovieIdAndTheaterId(ctx context.Context, movieId, theaterId int) ([]MovieSchedule, error)
	GetMovieScheduleByMovieIdAndShowTimeId(ctx context.Context, movieId, showTimeId int) ([]MovieSchedule, error)
	GetMovieScheduleByTheaterIdAndShowTimeId(ctx context.Context, theaterId, showTimeId int) ([]MovieSchedule, error)
	GetMovieScheduleByID(ctx context.Context, id int) (*MovieSchedule, error)
	DeleteMovieScheduleById(ctx context.Context, id int) error
	DeleteMovieScheduleByMovieIdAndTheaterId(ctx context.Context, movieId, theaterId int) error
	DeleteMovieScheduleByMovieIdAndTheaterIdAndShowTimeId(ctx context.Context, movieId, theaterId, showTimeId int) error
}

type service struct {
	userAdmin user_admin.AdminServiceClient
	auth      auth.JWT_TokenServiceClient
}

func NewService(pb user_admin.AdminServiceClient, auth auth.JWT_TokenServiceClient) Service {
	return &service{
		userAdmin: pb,
		auth:      auth,
	}
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

func (s *service) Register(ctx context.Context, signUpData *Admin) error {
	reqData := user_admin.RegisterAdminRequest{
		Username:  signUpData.Username,
		Password:  signUpData.Password,
		Phone:     signUpData.PhoneNumber,
		Email:     signUpData.Email,
		FirstName: signUpData.FirstName,
		LastName:  signUpData.LastName,
		Gender:    signUpData.Gender,
	}
	if _, err := s.userAdmin.RegisterAdmin(ctx, &reqData); err != nil {
		return err
	}
	return nil
}
func (s *service) Login(ctx context.Context, loginData *Admin) (string, error) {
	req := user_admin.LoginAdminRequest{
		Email:    loginData.Email,
		Password: loginData.Password,
	}
	res, err := s.userAdmin.LoginAdmin(ctx, &req)
	if err != nil {
		return "", err
	}
	return res.Token, nil
}

func (s *service) GetAdminProfile(ctx context.Context, id int) (*Admin, error) {
	admin, err := s.userAdmin.GetAdminProfile(ctx, &user_admin.GetProfileRequest{
		UserId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	return &Admin{
		ID:          int(admin.ProfileDetails.Id),
		Username:    admin.ProfileDetails.Username,
		PhoneNumber: admin.ProfileDetails.Phone,
		Email:       admin.ProfileDetails.Email,
		FirstName:   admin.ProfileDetails.FirstName,
		LastName:    admin.ProfileDetails.LastName,
		DateOfBirth: admin.ProfileDetails.DateOfBirth,
		Gender:      admin.ProfileDetails.Gender,
		IsVerified:  admin.ProfileDetails.IsVerified,
	}, nil
}
func (s *service) UpdateAdminProfile(ctx context.Context, id int, admin AdminProfileDetails) error {
	_, err := s.userAdmin.UpdateAdminProfile(ctx, &user_admin.UpdateAdminProfileRequest{
		UserId:      int32(id),
		Username:    admin.Username,
		Phone:       admin.PhoneNumber,
		FirstName:   admin.FirstName,
		LastName:    admin.LastName,
		Gender:      admin.Gender,
		DateOfBirth: admin.DateOfBirth,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ForgotPassword(ctx context.Context, email string) error {
	_, err := s.userAdmin.ForgotAdminPassword(ctx, &user_admin.ForgotPasswordRequest{
		Email: email,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ResetPassword(ctx context.Context, data ResetPassword) error {
	_, err := s.userAdmin.ResetAdminPassword(ctx, &user_admin.ResetPasswordRequest{
		Email:       data.Email,
		Otp:         data.Otp,
		NewPassword: data.NewPassword,
	})
	if err != nil {
		return err
	}
	return nil
}

// Theater
func (s *service) AddTheater(ctx context.Context, theater Theater) error {
	_, err := s.userAdmin.AddTheater(ctx, &user_admin.AddTheaterRequest{
		Name:            theater.Name,
		Place:           theater.Place,
		City:            theater.City,
		District:        theater.District,
		State:           theater.State,
		OwnerId:         uint32(theater.OwnerID),
		NumberOfScreens: int32(theater.NumberOfScreens),
		TheaterTypeId:   int32(theater.TheaterTypeID),
	},
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteTheaterByID(ctx context.Context, id int) error {
	_, err := s.userAdmin.DeleteTheaterByID(ctx, &user_admin.DeleteTheaterRequest{
		TheaterId: int32(id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteTheaterByName(ctx context.Context, name string) error {
	_, err := s.userAdmin.DeleteTheaterByName(ctx, &user_admin.DeleteTheaterByNameRequest{
		Name: name,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetTheaterByID(ctx context.Context, id int) (*Theater, error) {
	response, err := s.userAdmin.GetTheaterByID(ctx, &user_admin.GetTheaterByIDRequest{
		TheaterId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	return &Theater{
		ID:              uint(response.Theater.TheaterId),
		Name:            response.Theater.Name,
		Place:           response.Theater.Place,
		City:            response.Theater.City,
		District:        response.Theater.District,
		State:           response.Theater.State,
		OwnerID:         uint(response.Theater.OwnerId),
		NumberOfScreens: int(response.Theater.NumberOfScreens),
		TheaterTypeID:   int(response.Theater.TheaterTypeId),
	}, nil
}

func (s *service) GetTheaterByName(ctx context.Context, name string) ([]Theater, error) {
	response, err := s.userAdmin.GetTheaterByName(ctx, &user_admin.GetTheaterByNameRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	theaters := []Theater{}
	for _, res := range response.Theater {
		theater := Theater{
			ID:              uint(res.TheaterId),
			Name:            res.Name,
			Place:           res.Place,
			City:            res.City,
			District:        res.District,
			State:           res.State,
			OwnerID:         uint(res.OwnerId),
			NumberOfScreens: int(res.NumberOfScreens),
			TheaterTypeID:   int(res.TheaterTypeId),
		}
		theaters = append(theaters, theater)
	}
	return theaters, nil
}

func (s *service) ListTheaters(ctx context.Context) ([]Theater, error) {
	response, err := s.userAdmin.ListTheaters(ctx, &user_admin.ListTheatersRequest{})
	if err != nil {
		return nil, err
	}
	theaters := []Theater{}

	for _, res := range response.Theaters {
		theater := Theater{
			ID:              uint(res.TheaterId),
			Name:            res.Name,
			Place:           res.Place,
			City:            res.City,
			District:        res.District,
			State:           res.State,
			OwnerID:         uint(res.OwnerId),
			NumberOfScreens: int(res.NumberOfScreens),
			TheaterTypeID:   int(res.TheaterTypeId),
		}
		theaters = append(theaters, theater)
	}
	return theaters, nil
}

func (s *service) UpdateTheater(ctx context.Context, id int, theater Theater) error {
	_, err := s.userAdmin.UpdateTheater(ctx, &user_admin.UpdateTheaterRequest{
		TheaterId:       int32(id),
		Name:            theater.Name,
		Place:           theater.Place,
		City:            theater.City,
		District:        theater.District,
		State:           theater.State,
		OwnerId:         uint32(theater.OwnerID),
		NumberOfScreens: int32(theater.NumberOfScreens),
		TheaterTypeId:   int32(theater.TheaterTypeID),
	},
	)
	if err != nil {
		return err
	}
	return nil
}

//Movies

func (s *service) ListMovies(ctx context.Context) ([]Movie, error) {
	response, err := s.userAdmin.ListMovies(ctx, &user_admin.ListMoviesRequest{})
	if err != nil {
		return nil, err
	}
	var movies []Movie
	for _, m := range response.Movies {
		movie := Movie{
			Title:       m.Title,
			Description: m.Description,
			Duration:    int(m.Duration),
			Genre:       m.Genre,
			ReleaseDate: m.ReleaseDate,
			Rating:      float64(m.Rating),
			Language:    m.Language,
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// Theater types
func (s *service) ListTheaterTypes(ctx context.Context) ([]TheaterType, error) {
	response, err := s.userAdmin.ListTheaterTypes(ctx, &user_admin.ListTheaterTypesRequest{})
	if err != nil {
		return nil, err
	}
	theaterTypes := []TheaterType{}

	for _, res := range response.TheaterTypes {
		theaterType := TheaterType{
			ID:              int(res.Id),
			TheaterTypeName: res.TheaterTypeName,
		}
		theaterTypes = append(theaterTypes, theaterType)
	}
	return theaterTypes, nil
}

//Screen types

func (s *service) ListScreenTypes(ctx context.Context) ([]ScreenType, error) {
	response, err := s.userAdmin.ListScreenTypes(ctx, &user_admin.ListScreenTypesRequest{})
	if err != nil {
		return nil, err
	}
	screenTypes := []ScreenType{}

	for _, res := range response.ScreenTypes {
		screenType := ScreenType{
			ID:             int(res.Id),
			ScreenTypeName: res.ScreenTypeName,
		}
		screenTypes = append(screenTypes, screenType)
	}
	return screenTypes, nil
}

//Seat categories

func (s *service) ListSeatCategories(ctx context.Context) ([]SeatCategory, error) {
	response, err := s.userAdmin.ListSeatCategories(ctx, &user_admin.ListSeatCategoriesRequest{})
	if err != nil {
		return nil, err
	}
	seatCategories := []SeatCategory{}

	for _, res := range response.SeatCategories {
		seatCategory := SeatCategory{
			ID:               int(res.Id),
			SeatCategoryName: res.SeatCategoryName,
		}
		seatCategories = append(seatCategories, seatCategory)
	}
	return seatCategories, nil
}

// Theater screen
// TheaterScreen
func (s *service) AddTheaterScreen(ctx context.Context, theaterScreen TheaterScreen) error {
	_, err := s.userAdmin.AddTheaterScreen(ctx, &user_admin.AddTheaterScreenRequest{
		TheaterScreen: &user_admin.TheaterScreen{
			ID:           uint32(theaterScreen.ID),
			TheaterID:    int32(theaterScreen.TheaterID),
			ScreenNumber: int32(theaterScreen.ScreenNumber),
			SeatCapacity: int32(theaterScreen.SeatCapacity),
			ScreenTypeID: int32(theaterScreen.ScreenTypeID),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteTheaterScreenByID(ctx context.Context, id int) error {
	_, err := s.userAdmin.DeleteTheaterScreenByID(ctx, &user_admin.DeleteTheaterScreenRequest{
		TheaterScreenId: int32(id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteTheaterScreenByNumber(ctx context.Context, theaterID, screenNumber int) error {
	_, err := s.userAdmin.DeleteTheaterScreenByNumber(ctx, &user_admin.DeleteTheaterScreenByNumberRequest{
		TheaterID:    int32(theaterID),
		ScreenNumber: int32(screenNumber),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetTheaterScreenByID(ctx context.Context, id int) (*TheaterScreen, error) {
	response, err := s.userAdmin.GetTheaterScreenByID(ctx, &user_admin.GetTheaterScreenByIDRequest{
		TheaterScreenId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	return &TheaterScreen{
		ID:           uint(response.TheaterScreen.ID),
		TheaterID:    int(response.TheaterScreen.TheaterID),
		ScreenNumber: int(response.TheaterScreen.ScreenNumber),
		SeatCapacity: int(response.TheaterScreen.SeatCapacity),
		ScreenTypeID: int(response.TheaterScreen.ScreenTypeID),
	}, nil
}

func (s *service) GetTheaterScreenByNumber(ctx context.Context, theaterID, screenNumber int) (*TheaterScreen, error) {
	response, err := s.userAdmin.GetTheaterScreenByNumber(ctx, &user_admin.GetTheaterScreenByNumberRequest{
		TheaterID:    int32(theaterID),
		ScreenNumber: int32(screenNumber),
	})
	if err != nil {
		return nil, err
	}
	return &TheaterScreen{
		ID:           uint(response.TheaterScreen.ID),
		TheaterID:    int(response.TheaterScreen.TheaterID),
		ScreenNumber: int(response.TheaterScreen.ScreenNumber),
		SeatCapacity: int(response.TheaterScreen.SeatCapacity),
		ScreenTypeID: int(response.TheaterScreen.ScreenTypeID),
	}, nil
}

func (s *service) ListTheaterScreens(ctx context.Context, theaterID int) ([]TheaterScreen, error) {
	response, err := s.userAdmin.ListTheaterScreens(ctx, &user_admin.ListTheaterScreensRequest{
		TheaterID: int32(theaterID),
	})
	if err != nil {
		return nil, err
	}
	theaterScreens := []TheaterScreen{}

	for _, res := range response.TheaterScreens {
		theaterScreen := TheaterScreen{
			ID:           uint(res.ID),
			TheaterID:    int(res.TheaterID),
			ScreenNumber: int(res.ScreenNumber),
			SeatCapacity: int(res.SeatCapacity),
			ScreenTypeID: int(res.ScreenTypeID),
		}
		theaterScreens = append(theaterScreens, theaterScreen)
	}
	return theaterScreens, nil
}

func (s *service) UpdateTheaterScreen(ctx context.Context, id int, theaterScreen TheaterScreen) error {
	_, err := s.userAdmin.UpdateTheaterScreen(ctx, &user_admin.UpdateTheaterScreenRequest{
		TheaterScreen: &user_admin.TheaterScreen{
			ID:           uint32(id),
			TheaterID:    int32(theaterScreen.TheaterID),
			ScreenNumber: int32(theaterScreen.ScreenNumber),
			SeatCapacity: int32(theaterScreen.SeatCapacity),
			ScreenTypeID: int32(theaterScreen.ScreenTypeID),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// Show Times
func (s *service) AddShowtime(ctx context.Context, showtime Showtime) error {
	_, err := s.userAdmin.AddShowtime(ctx, &user_admin.AddShowtimeRequest{
		Showtime: &user_admin.Showtime{
			Id:       uint32(showtime.ID),
			MovieId:  int32(showtime.MovieID),
			ScreenId: int32(showtime.ScreenID),
			ShowDate: timestamppb.New(showtime.ShowDate),
			ShowTime: timestamppb.New(showtime.ShowTime),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteShowtimeByID(ctx context.Context, id int) error {
	_, err := s.userAdmin.DeleteShowtimeByID(ctx, &user_admin.DeleteShowtimeRequest{
		ShowtimeId: int32(id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteShowtimeByDetails(ctx context.Context, movieID, screenID int, showDate, showTime time.Time) error {
	_, err := s.userAdmin.DeleteShowtimeByDetails(ctx, &user_admin.DeleteShowtimeByDetailsRequest{
		MovieId:  int32(movieID),
		ScreenId: int32(screenID),
		ShowDate: timestamppb.New(showDate),
		ShowTime: timestamppb.New(showTime),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetShowtimeByID(ctx context.Context, id int) (*Showtime, error) {
	response, err := s.userAdmin.GetShowtimeByID(ctx, &user_admin.GetShowtimeByIDRequest{
		ShowtimeId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	return &Showtime{
		ID:       uint(response.Showtime.Id),
		MovieID:  int(response.Showtime.MovieId),
		ScreenID: int(response.Showtime.ScreenId),
		ShowDate: response.Showtime.ShowDate.AsTime(),
		ShowTime: response.Showtime.ShowDate.AsTime(),
	}, nil
}

func (s *service) GetShowtimeByDetails(ctx context.Context, movieID, screenID int, showDate, showTime time.Time) (*Showtime, error) {
	response, err := s.userAdmin.GetShowtimeByDetails(ctx, &user_admin.GetShowtimeByDetailsRequest{
		MovieId:  int32(movieID),
		ScreenId: int32(screenID),
		ShowDate: timestamppb.New(showDate),
		ShowTime: timestamppb.New(showTime),
	})
	if err != nil {
		return nil, err
	}
	return &Showtime{
		ID:       uint(response.Showtime.Id),
		MovieID:  int(response.Showtime.MovieId),
		ScreenID: int(response.Showtime.ScreenId),
		ShowDate: response.Showtime.ShowDate.AsTime(),
		ShowTime: response.Showtime.ShowDate.AsTime(),
	}, nil
}

func (s *service) UpdateShowtime(ctx context.Context, id int, showtime Showtime) error {
	_, err := s.userAdmin.UpdateShowtime(ctx, &user_admin.UpdateShowtimeRequest{
		Showtime: &user_admin.Showtime{
			Id:       uint32(id),
			MovieId:  int32(showtime.MovieID),
			ScreenId: int32(showtime.ScreenID),
			ShowDate: timestamppb.New(showtime.ShowDate),
			ShowTime: timestamppb.New(showtime.ShowTime),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ListShowtimes(ctx context.Context, movieID int) ([]Showtime, error) {
	response, err := s.userAdmin.ListShowtimes(ctx, &user_admin.ListShowtimesRequest{
		MovieId: int32(movieID),
	})
	if err != nil {
		return nil, err
	}
	showtimes := []Showtime{}

	for _, res := range response.Showtimes {
		showtime := Showtime{
			ID:       uint(res.Id),
			MovieID:  int(res.MovieId),
			ScreenID: int(res.ScreenId),
			ShowDate: res.ShowDate.AsTime(),
			ShowTime: res.ShowTime.AsTime(), // Convert string to time.Time if needed
		}
		showtimes = append(showtimes, showtime)
	}
	return showtimes, nil
}

// Movie Schedule
func (s *service) AddMovieSchedule(ctx context.Context, movieSchedule MovieSchedule) error {
	_, err := s.userAdmin.AddMovieSchedule(ctx, &user_admin.AddMovieScheduleRequest{
		MovieSchedule: &user_admin.MovieSchedule{
			MovieId:    int32(movieSchedule.MovieID),
			TheaterId:  int32(movieSchedule.TheaterID),
			ShowtimeId: int32(movieSchedule.ShowtimeID),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteMovieScheduleById(ctx context.Context, id int) error {
	_, err := s.userAdmin.DeleteMovieScheduleById(ctx, &user_admin.DeleteMovieScheduleByIdRequest{
		Id: int32(id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteMovieScheduleByMovieIdAndTheaterId(ctx context.Context, movieId int, theaterId int) error {
	_, err := s.userAdmin.DeleteMovieScheduleByMovieIdAndTheaterId(ctx, &user_admin.DeleteMovieScheduleByMovieIdAndTheaterIdRequest{
		MovieId:   int32(movieId),
		TheaterId: int32(theaterId),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteMovieScheduleByMovieIdAndTheaterIdAndShowTimeId(ctx context.Context, movieId int, theaterId int, showTimeId int) error {
	_, err := s.userAdmin.DeleteMovieScheduleByMovieIdAndTheaterIdAndShowTimeId(ctx, &user_admin.DeleteMovieScheduleByMovieIdAndTheaterIdAndShowTimeIdRequest{
		MovieId:    int32(movieId),
		TheaterId:  int32(theaterId),
		ShowtimeId: int32(showTimeId),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAllMovieSchedules(ctx context.Context) ([]MovieSchedule, error) {
	response, err := s.userAdmin.GetAllMovieSchedules(ctx, &user_admin.GetAllMovieScheduleRequest{})
	if err != nil {
		return nil, err
	}
	movieSchedules := []MovieSchedule{}

	for _, res := range response.MovieSchedules {
		movieSchedule := MovieSchedule{
			ID:         uint(res.Id),
			MovieID:    int(res.MovieId),
			TheaterID:  int(res.TheaterId),
			ShowtimeID: int(res.ShowtimeId),
		}
		movieSchedules = append(movieSchedules, movieSchedule)
	}
	return movieSchedules, nil
}

func (s *service) GetMovieScheduleByID(ctx context.Context, id int) (*MovieSchedule, error) {
	response, err := s.userAdmin.GetMovieScheduleByID(ctx, &user_admin.GetMovieScheduleByIDRequest{
		Id: int32(id),
	})

	if err != nil {
		return nil, err
	}
	return &MovieSchedule{
		ID:         uint(response.MovieSchedule.Id),
		MovieID:    int(response.MovieSchedule.MovieId),
		TheaterID:  int(response.MovieSchedule.TheaterId),
		ShowtimeID: int(response.MovieSchedule.ShowtimeId),
	}, nil
}

func (s *service) GetMovieScheduleByMovieID(ctx context.Context, movieId int) ([]MovieSchedule, error) {
	response, err := s.userAdmin.GetMovieScheduleByMovieID(ctx, &user_admin.GetMovieScheduleByMovieIdRequest{
		MovieId: int32(movieId),
	})
	if err != nil {
		return nil, err
	}
	var movieSchedules []MovieSchedule
	for _, res := range response.MovieSchedules {
		movieSchedule := MovieSchedule{
			ID:         uint(res.Id),
			MovieID:    int(res.MovieId),
			TheaterID:  int(res.TheaterId),
			ShowtimeID: int(res.ShowtimeId),
		}
		movieSchedules = append(movieSchedules, movieSchedule)
	}
	return movieSchedules, nil
}

func (s *service) GetMovieScheduleByMovieIdAndShowTimeId(ctx context.Context, movieId int, showTimeId int) ([]MovieSchedule, error) {
	response, err := s.userAdmin.GetMovieScheduleByMovieIdAndShowTimeId(ctx, &user_admin.GetMovieScheduleByMovieIdAndShowTimeIdRequest{
		MovieId:    int32(movieId),
		ShowtimeId: int32(showTimeId),
	})
	if err != nil {
		return nil, err
	}
	var movieSchedules []MovieSchedule
	for _, res := range response.MovieSchedules {
		movieSchedule := MovieSchedule{
			ID:         uint(res.Id),
			MovieID:    int(res.MovieId),
			TheaterID:  int(res.TheaterId),
			ShowtimeID: int(res.ShowtimeId),
		}
		movieSchedules = append(movieSchedules, movieSchedule)
	}
	return movieSchedules, nil
}

func (s *service) GetMovieScheduleByMovieIdAndTheaterId(ctx context.Context, movieId int, theaterId int) ([]MovieSchedule, error) {
	response, err := s.userAdmin.GetMovieScheduleByMovieIdAndTheaterId(ctx, &user_admin.GetMovieScheduleByMovieIdAndTheaterIdRequest{
		MovieId:   int32(movieId),
		TheaterId: int32(theaterId),
	})
	if err != nil {
		return nil, err
	}
	var movieSchedules []MovieSchedule
	for _, res := range response.MovieSchedules {
		movieSchedule := MovieSchedule{
			ID:         uint(res.Id),
			MovieID:    int(res.MovieId),
			TheaterID:  int(res.TheaterId),
			ShowtimeID: int(res.ShowtimeId),
		}
		movieSchedules = append(movieSchedules, movieSchedule)
	}
	return movieSchedules, nil
}

func (s *service) GetMovieScheduleByTheaterID(ctx context.Context, theaterId int) ([]MovieSchedule, error) {
	response, err := s.userAdmin.GetMovieScheduleByTheaterID(ctx, &user_admin.GetMovieScheduleByTheaterIdRequest{
		TheaterId: int32(theaterId),
	})
	if err != nil {
		return nil, err
	}
	var movieSchedules []MovieSchedule
	for _, res := range response.MovieSchedules {
		movieSchedule := MovieSchedule{
			ID:         uint(res.Id),
			MovieID:    int(res.MovieId),
			TheaterID:  int(res.TheaterId),
			ShowtimeID: int(res.ShowtimeId),
		}
		movieSchedules = append(movieSchedules, movieSchedule)
	}
	return movieSchedules, nil
}

func (s *service) GetMovieScheduleByTheaterIdAndShowTimeId(ctx context.Context, theaterId int, showTimeId int) ([]MovieSchedule, error) {
	response, err := s.userAdmin.GetMovieScheduleByTheaterIdAndShowTimeId(ctx, &user_admin.GetGetMovieScheduleByTheaterIdAndShowTimeIdRequest{
		TheaterId:  int32(theaterId),
		ShowtimeId: int32(showTimeId),
	})
	if err != nil {
		return nil, err
	}
	var movieSchedules []MovieSchedule
	for _, res := range response.MovieSchedules {
		movieSchedule := MovieSchedule{
			ID:         uint(res.Id),
			MovieID:    int(res.MovieId),
			TheaterID:  int(res.TheaterId),
			ShowtimeID: int(res.ShowtimeId),
		}
		movieSchedules = append(movieSchedules, movieSchedule)
	}
	return movieSchedules, nil
}

func (s *service) UpdateMovieSchedule(ctx context.Context, id int, updateData MovieSchedule) error {
	_, err := s.userAdmin.UpdateMovieSchedule(ctx, &user_admin.UpdateMovieScheduleRequest{
		MovieSchedule: &user_admin.MovieSchedule{
			Id:         int32(id),
			MovieId:    int32(updateData.MovieID),
			TheaterId:  int32(updateData.TheaterID),
			ShowtimeId: int32(updateData.ShowtimeID),
		},
	})
	if err != nil {
		return err
	}
	return nil
}
