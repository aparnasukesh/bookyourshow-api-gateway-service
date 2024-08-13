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
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

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
