package admin

import (
	"context"
	"strings"

	"github.com/aparnasukesh/inter-communication/auth"
	"github.com/aparnasukesh/inter-communication/user_admin"
)

type Service interface {
	Register(ctx context.Context, signUpData *Admin) error
	Login(ctx context.Context, loginData *Admin) (string, error)
	GetUserIDFromToken(ctx context.Context, authorization string) (int, error)
	//Theater
	AddTheater(ctx context.Context, theater Theater) error
	DeleteTheaterByID(ctx context.Context, id int) error
	DeleteTheaterByName(ctx context.Context, name string) error
	GetTheaterByID(ctx context.Context, id int) (*Theater, error)
	GetTheaterByName(ctx context.Context, name string) (*Theater, error)
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

// Theater
func (s *service) AddTheater(ctx context.Context, theater Theater) error {
	_, err := s.userAdmin.AddTheater(ctx, &user_admin.AddTheaterRequest{
		Name:            theater.Name,
		Location:        theater.Location,
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
		ID:              int(response.Theater.TheaterId),
		Name:            response.Theater.Name,
		Location:        response.Theater.Location,
		OwnerID:         uint(response.Theater.OwnerId),
		NumberOfScreens: int(response.Theater.NumberOfScreens),
		TheaterTypeID:   int(response.Theater.TheaterTypeId),
	}, nil
}

func (s *service) GetTheaterByName(ctx context.Context, name string) (*Theater, error) {
	response, err := s.userAdmin.GetTheaterByName(ctx, &user_admin.GetTheaterByNameRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return &Theater{
		ID:              int(response.Theater.TheaterId),
		Name:            response.Theater.Name,
		Location:        response.Theater.Location,
		OwnerID:         uint(response.Theater.OwnerId),
		NumberOfScreens: int(response.Theater.NumberOfScreens),
		TheaterTypeID:   int(response.Theater.TheaterTypeId),
	}, nil
}

func (s *service) ListTheaters(ctx context.Context) ([]Theater, error) {
	response, err := s.userAdmin.ListTheaters(ctx, &user_admin.ListTheatersRequest{})
	if err != nil {
		return nil, err
	}
	theaters := []Theater{}

	for _, res := range response.Theaters {
		theater := Theater{
			ID:              int(res.TheaterId),
			Name:            res.Name,
			Location:        res.Location,
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
		Location:        theater.Location,
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
			ID:                int(res.Id),
			SeatCategoryName:  res.SeatCategoryName,
			SeatCategoryPrice: res.SeatCategoryPrice,
		}
		seatCategories = append(seatCategories, seatCategory)
	}
	return seatCategories, nil
}
