package superadmin

import (
	"context"
	"errors"

	"github.com/aparnasukesh/inter-communication/auth"
	"github.com/aparnasukesh/inter-communication/movie_booking"
	"github.com/aparnasukesh/inter-communication/user_admin"
)

type service struct {
	userAdmin    user_admin.SuperAdminServiceClient
	movieBooking movie_booking.MovieServiceClient
}
type Service interface {
	Login(ctx context.Context, loginData *Admin) (string, error)
	ListAdminRequests(ctx context.Context) ([]AdminRequestResponse, error)
	AdminApproval(ctx context.Context, email string, isVerified bool) error
	ListAllAdmins(ctx context.Context) ([]Admin, error)
	GetAdminById(ctx context.Context, id int) (*Admin, error)
	// User
	ListAllUser(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	BlockUser(ctx context.Context, id int) error
	UnBlockUser(ctx context.Context, id int) error
	// Movies
	RegisterMovie(ctx context.Context, movie Movie) (int, error)
	UpdateMovie(ctx context.Context, movie Movie, movieId int) error
	ListMovies(ctx context.Context) ([]Movie, error)
	GetMovieDetails(ctx context.Context, movieId int) (*Movie, error)
	DeleteMovie(ctx context.Context, movieId int) error
	// Theater-Type
	AddTheaterType(ctx context.Context, data TheaterType) error
	DeleteTheaterTypeById(ctx context.Context, id int) error
	DeleteTheaterTypeByName(ctx context.Context, theaterName string) error
	GetTheaterTypeByID(ctx context.Context, id int) (*TheaterType, error)
	GetTheaterTypeByName(ctx context.Context, name string) (*TheaterType, error)
	UpdateTheaterType(ctx context.Context, id int, theaterType TheaterType) error
	ListTheaterTypes(ctx context.Context) ([]TheaterType, error)
	// Screen-Type
	AddScreenType(ctx context.Context, data ScreenType) error
	DeleteScreenTypeById(ctx context.Context, id int) error
	DeleteScreenTypeByName(ctx context.Context, screenName string) error
	GetScreenTypeByID(ctx context.Context, id int) (*ScreenType, error)
	GetScreenTypeByName(ctx context.Context, name string) (*ScreenType, error)
	UpdateScreenType(ctx context.Context, id int, screenType ScreenType) error
	ListScreenTypes(ctx context.Context) ([]ScreenType, error)
	// Seat category
	AddSeatCategory(ctx context.Context, seatCategory SeatCategory) error
	DeleteSeatCategoryByID(ctx context.Context, id int) error
	DeleteSeatCategoryByName(ctx context.Context, name string) error
	GetSeatCategoryByID(ctx context.Context, id int) (*SeatCategory, error)
	GetSeatCategoryByName(ctx context.Context, name string) (*SeatCategory, error)
	UpdateSeatCategory(ctx context.Context, id int, seatCategory SeatCategory) error
	ListSeatCategories(ctx context.Context) ([]SeatCategory, error)
}

func NewService(pb user_admin.SuperAdminServiceClient, auth auth.JWT_TokenServiceClient, movieBooking movie_booking.MovieServiceClient) Service {
	return &service{
		userAdmin:    pb,
		movieBooking: movieBooking,
	}
}
func (s *service) Login(ctx context.Context, loginData *Admin) (string, error) {
	req := user_admin.LoginSuperAdminRequest{
		Email:    loginData.Email,
		Password: loginData.Password,
	}
	res, err := s.userAdmin.LoginSuperAdmin(ctx, &req)
	if err != nil {
		return "", err
	}
	return res.Token, nil
}
func (s *service) ListAdminRequests(ctx context.Context) ([]AdminRequestResponse, error) {
	res, err := s.userAdmin.ListAdminRequests(ctx, &user_admin.ListAdminRequestsRequest{})
	if err != nil {
		return nil, err
	}

	adminRequests := make([]AdminRequestResponse, len(res.Email))
	for i, admin := range res.Email {
		adminRequests[i] = AdminRequestResponse{
			Email: admin.Email,
		}
	}
	return adminRequests, nil
}

func (s *service) AdminApproval(ctx context.Context, email string, isVerified bool) error {
	_, err := s.userAdmin.AdminApproval(ctx, &user_admin.AdminApprovalRequest{
		Email:      email,
		IsVerified: isVerified,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ListAllAdmins(ctx context.Context) ([]Admin, error) {
	response, err := s.userAdmin.ListAllAdmin(ctx, &user_admin.ListAllAdminRequest{})
	if err != nil {
		return nil, err
	}
	admins := []Admin{}

	for _, res := range response.Admin {
		admin := Admin{
			ID:          int(res.Id),
			Username:    res.Username,
			PhoneNumber: res.Phone,
			Email:       res.Email,
			FirstName:   res.FirstName,
			LastName:    res.LastName,
			Gender:      res.Gender,
		}
		admins = append(admins, admin)
	}
	return admins, nil

}

func (s *service) GetAdminById(ctx context.Context, id int) (*Admin, error) {
	response, err := s.userAdmin.GetAdminByID(ctx, &user_admin.GetAdminByIdRequest{
		AdminId: int32(id),
	})

	if err != nil {
		return nil, err
	}
	return &Admin{
		ID:          id,
		Username:    response.Admin.Username,
		PhoneNumber: response.Admin.Phone,
		Email:       response.Admin.Email,
		FirstName:   response.Admin.FirstName,
		LastName:    response.Admin.LastName,
		Gender:      response.Admin.Gender,
	}, nil
}

// Movies
func (s *service) RegisterMovie(ctx context.Context, movie Movie) (int, error) {
	response, err := s.userAdmin.RegisterMovie(ctx, &user_admin.RegisterMovieRequest{
		Title:       movie.Title,
		Description: movie.Description,
		Duration:    int32(movie.Duration),
		Genre:       movie.Genre,
		ReleaseDate: movie.ReleaseDate,
		Rating:      float32(movie.Rating),
		Language:    movie.Language,
	})
	if err != nil {
		return 0, err
	}
	return int(response.MovieId), nil
}

func (s *service) UpdateMovie(ctx context.Context, movie Movie, movieId int) error {
	_, err := s.userAdmin.UpdateMovie(ctx, &user_admin.UpdateMovieRequest{
		MovieId:     uint32(movieId),
		Title:       movie.Title,
		Description: movie.Description,
		Duration:    int32(movie.Duration),
		Genre:       movie.Genre,
		ReleaseDate: movie.ReleaseDate,
		Rating:      float32(movie.Rating),
		Language:    movie.Language,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteMovie(ctx context.Context, movieId int) error {
	_, err := s.userAdmin.DeleteMovie(ctx, &user_admin.DeleteMovieRequest{
		MovieId: uint32(movieId),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetMovieDetails(ctx context.Context, movieId int) (*Movie, error) {
	response, err := s.userAdmin.GetMovieDetails(ctx, &user_admin.GetMovieDetailsRequest{
		MovieId: uint32(movieId),
	})
	if err != nil {
		return nil, err
	}
	if response.Movie == nil {
		return nil, errors.New("movie details not found")
	}
	movie := &Movie{
		Title:       response.Movie.Title,
		Description: response.Movie.Description,
		Duration:    int(response.Movie.Duration),
		Genre:       response.Movie.Genre,
		ReleaseDate: response.Movie.ReleaseDate,
		Rating:      float64(response.Movie.Rating),
		Language:    response.Movie.Language,
	}

	return movie, nil
}

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

// Theater type
func (s *service) AddTheaterType(ctx context.Context, data TheaterType) error {
	_, err := s.userAdmin.AddTheaterType(ctx, &user_admin.AddTheaterTypeRequest{
		TheaterTypeName: data.TheaterTypeName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteTheaterTypeById(ctx context.Context, id int) error {
	_, err := s.userAdmin.DeleteTheaterTypeByID(ctx, &user_admin.DeleteTheaterTypeRequest{
		TheaterTypeId: int32(id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteTheaterTypeByName(ctx context.Context, theaterName string) error {
	_, err := s.userAdmin.DeleteTheaterTypeByName(ctx, &user_admin.DeleteTheaterTypeByNameRequest{
		Name: theaterName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetTheaterTypeByID(ctx context.Context, id int) (*TheaterType, error) {
	response, err := s.userAdmin.GetTheaterTypeByID(ctx, &user_admin.GetTheaterTypeByIDRequest{
		TheaterTypeId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	return &TheaterType{
		ID:              id,
		TheaterTypeName: response.TheaterType.TheaterTypeName,
	}, nil
}

func (s *service) GetTheaterTypeByName(ctx context.Context, name string) (*TheaterType, error) {
	response, err := s.userAdmin.GetTheaterTypeByName(ctx, &user_admin.GetTheaterTypeByNameRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return &TheaterType{
		ID:              int(response.TheaterType.Id),
		TheaterTypeName: name,
	}, nil
}

func (s *service) UpdateTheaterType(ctx context.Context, id int, theaterType TheaterType) error {
	_, err := s.userAdmin.UpdateTheaterType(ctx, &user_admin.UpdateTheaterTypeRequest{
		Id:              int32(id),
		TheaterTypeName: theaterType.TheaterTypeName,
	})
	if err != nil {
		return err
	}
	return nil
}

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

// Screen type
func (s *service) AddScreenType(ctx context.Context, data ScreenType) error {
	_, err := s.userAdmin.AddScreenType(ctx, &user_admin.AddScreenTypeRequest{
		ScreenTypeName: data.ScreenTypeName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteScreenTypeById(ctx context.Context, id int) error {
	_, err := s.userAdmin.DeleteScreenTypeByID(ctx, &user_admin.DeleteScreenTypeRequest{
		ScreenTypeId: int32(id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteScreenTypeByName(ctx context.Context, screenName string) error {
	_, err := s.userAdmin.DeleteScreenTypeByName(ctx, &user_admin.DeleteScreenTypeByNameRequest{
		Name: screenName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetScreenTypeByID(ctx context.Context, id int) (*ScreenType, error) {
	response, err := s.userAdmin.GetScreenTypeByID(ctx, &user_admin.GetScreenTypeByIDRequest{
		ScreenTypeId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	return &ScreenType{
		ID:             id,
		ScreenTypeName: response.ScreenType.ScreenTypeName,
	}, nil
}

func (s *service) GetScreenTypeByName(ctx context.Context, name string) (*ScreenType, error) {
	response, err := s.userAdmin.GetScreenTypeByName(ctx, &user_admin.GetScreenTypeByNameRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return &ScreenType{
		ID:             int(response.ScreenType.Id),
		ScreenTypeName: name,
	}, nil
}

func (s *service) UpdateScreenType(ctx context.Context, id int, screenType ScreenType) error {
	_, err := s.userAdmin.UpdateScreenType(ctx, &user_admin.UpdateScreenTypeRequest{
		Id:             int32(id),
		ScreenTypeName: screenType.ScreenTypeName,
	})
	if err != nil {
		return err
	}
	return nil
}

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

// seat category
func (s *service) AddSeatCategory(ctx context.Context, seatCategory SeatCategory) error {
	_, err := s.userAdmin.AddSeatCategory(ctx, &user_admin.AddSeatCategoryRequest{
		SeatCategory: &user_admin.SeatCategory{
			SeatCategoryName: seatCategory.SeatCategoryName,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteSeatCategoryByID(ctx context.Context, id int) error {
	_, err := s.userAdmin.DeleteSeatCategoryByID(ctx, &user_admin.DeleteSeatCategoryRequest{
		SeatCategoryId: int32(id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteSeatCategoryByName(ctx context.Context, name string) error {
	_, err := s.userAdmin.DeleteSeatCategoryByName(ctx, &user_admin.DeleteSeatCategoryByNameRequest{
		Name: name,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetSeatCategoryByID(ctx context.Context, id int) (*SeatCategory, error) {
	response, err := s.userAdmin.GetSeatCategoryByID(ctx, &user_admin.GetSeatCategoryByIDRequest{
		SeatCategoryId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	return &SeatCategory{
		ID:               int(response.SeatCategory.Id),
		SeatCategoryName: response.SeatCategory.SeatCategoryName,
	}, nil
}

func (s *service) GetSeatCategoryByName(ctx context.Context, name string) (*SeatCategory, error) {
	response, err := s.userAdmin.GetSeatCategoryByName(ctx, &user_admin.GetSeatCategoryByNameRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return &SeatCategory{
		ID:               int(response.SeatCategory.Id),
		SeatCategoryName: response.SeatCategory.SeatCategoryName,
	}, nil
}

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

func (s *service) UpdateSeatCategory(ctx context.Context, id int, seatCategory SeatCategory) error {
	_, err := s.userAdmin.UpdateSeatCategory(ctx, &user_admin.UpdateSeatCategoryRequest{
		Id: int32(id),
		SeatCategory: &user_admin.SeatCategory{
			Id:               int32(id),
			SeatCategoryName: seatCategory.SeatCategoryName,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// User
func (s *service) BlockUser(ctx context.Context, id int) error {
	_, err := s.userAdmin.BlockUser(ctx, &user_admin.BlockUserRequest{
		UserId: int32(id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetUserByID(ctx context.Context, id int) (*User, error) {
	response, err := s.userAdmin.GetUserByID(ctx, &user_admin.GetUserByIdRequest{
		UserId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	return &User{
		ID:          id,
		Username:    response.User.Username,
		Email:       response.User.Email,
		PhoneNumber: response.User.Phone,
		FirstName:   response.User.FirstName,
		LastName:    response.User.LastName,
		DateOfBirth: response.User.DateOfBirth,
		Gender:      response.User.Gender,
		IsVerified:  response.User.IsVerified,
	}, nil
}

func (s *service) ListAllUser(ctx context.Context) ([]User, error) {
	response, err := s.userAdmin.ListAllUser(ctx, &user_admin.ListAllUserRequest{})
	if err != nil {
		return nil, err
	}
	users := []User{}

	for _, res := range response.User {
		user := User{
			ID:          int(res.Id),
			Username:    res.Username,
			Email:       res.Email,
			PhoneNumber: res.Phone,
			FirstName:   res.FirstName,
			LastName:    res.LastName,
			DateOfBirth: res.DateOfBirth,
			Gender:      res.Gender,
			IsVerified:  res.IsVerified,
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *service) UnBlockUser(ctx context.Context, id int) error {
	_, err := s.userAdmin.UnBlockUser(ctx, &user_admin.UnBlockUserRequest{
		UserId: int32(id),
	})
	if err != nil {
		return err
	}
	return nil
}
