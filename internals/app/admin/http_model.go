package admin

import "time"

type Admin struct {
	ID          int    `json:"id"`
	Username    string `json:"username" validate:"required,min=8,max=24"`
	Password    string `json:"password" validate:"required,min=6,max=12"`
	PhoneNumber string `json:"phone" validate:"required,len=10"`
	Email       string `json:"email" validate:"email,required"`
	FirstName   string `gorm:"not null" json:"firstname" validate:"required,min=4,max=10"`
	LastName    string `gorm:"not null" json:"lastname" validate:"required,min=4,max=10"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	IsVerified  bool   `json:"is_verified"`
	OTP         string `json:"otp"`
}

type AdminProfileDetails struct {
	Username    string `json:"username" validate:"required,min=8,max=24"`
	FirstName   string `gorm:"not null" json:"firstname" validate:"required,min=4,max=10"`
	LastName    string `gorm:"not null" json:"lastname" validate:"required,min=4,max=10"`
	PhoneNumber string `json:"phone" validate:"required,len=10"`
	Email       string `json:"email" validate:"email,required"`
	DateOfBirth string `json:"dateofbirth"`
	Gender      string `json:"gender"`
}

// Movies
type Movie struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Duration    int     `json:"duration"`
	Genre       string  `json:"genre"`
	ReleaseDate string  `json:"release_date"`
	Rating      float64 `json:"rating"`
	Language    string  `json:"language"`
}

// Theater type
type TheaterType struct {
	ID              int    `json:"id"`
	TheaterTypeName string `json:"theater_type_name"`
}

type ScreenType struct {
	ID             int    `json:"id"`
	ScreenTypeName string `json:"screen_type_name"`
}

type SeatCategory struct {
	ID               int    `json:"id"`
	SeatCategoryName string `json:"seat_category_name"`
}

type Theater struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Place           string `json:"place"`
	City            string `json:"city"`
	District        string `json:"district"`
	State           string `json:"state"`
	OwnerID         uint   `json:"owner_id"`
	NumberOfScreens int    `json:"number_of_screens"`
	TheaterTypeID   int    `json:"theater_type_id"`
}

type TheaterScreen struct {
	ID           uint `json:"id"`
	TheaterID    int  `json:"theater_id"`
	ScreenNumber int  `json:"screen_number"`
	SeatCapacity int  `json:"seat_capacity"`
	ScreenTypeID int  `json:"screen_type_id"`
}

type Showtime struct {
	ID       uint      `json:"id"`
	MovieID  int       `json:"movie_id"`
	ScreenID int       `json:"screen_id"`
	ShowDate time.Time `json:"show_date"`
	ShowTime time.Time `json:"show_time"`
}

type ForgotPassword struct {
	Email string `json:"email"`
}

type ResetPassword struct {
	Email       string `json:"email"`
	Otp         string `json:"otp"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=12"`
}

type MovieSchedule struct {
	ID         uint `json:"id"`
	MovieID    int  `json:"movie_id"`
	TheaterID  int  `json:"theater_id"`
	ShowtimeID int  `json:"showtime_id"`
}
type RowSeatCategoryPrice struct {
	RowStart          string  `json:"row_start"`
	RowEnd            string  `json:"row_end"`
	SeatCategoryId    int     `json:"seat_category_id"`
	SeatCategoryPrice float32 `json:"seat_category_price"`
}

type CreateSeatsRequest struct {
	ID           int                    `json:"id"`
	ScreenId     int                    `json:"screen_id"`
	TotalRows    int                    `json:"total_rows"`
	TotalColumns int                    `json:"total_columns"`
	SeatRequest  []RowSeatCategoryPrice `json:"seat_request"`
}

type Seat struct {
	ID                int     `json:"id"`
	ScreenID          int     `json:"screen_id"`
	SeatNumber        string  `json:"seat_number"`
	Row               string  `json:"row"`
	Column            int     `json:"column"`
	SeatCategoryID    int     `json:"seat_category_id"`
	SeatCategoryPrice float64 `json:"seat_category_price"`
}
