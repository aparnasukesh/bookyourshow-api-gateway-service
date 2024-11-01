package user

import "time"

type User struct {
	Username    string `json:"username" validate:"required,min=8,max=24"`
	Password    string `json:"password" validate:"required,min=6,max=12"`
	PhoneNumber string `json:"phone" validate:"required,len=10"`
	Email       string `json:"email" validate:"email,required"`
	FirstName   string `gorm:"not null" json:"firstname" validate:"required,min=4,max=10"`
	LastName    string `gorm:"not null" json:"lastname" validate:"required,min=4,max=10"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	Otp         string `json:"otp"`
}

type UserProfileDetails struct {
	Username    string `json:"username" validate:"required,min=8,max=24"`
	FirstName   string `gorm:"not null" json:"firstname" validate:"required,min=4,max=10"`
	LastName    string `gorm:"not null" json:"lastname" validate:"required,min=4,max=10"`
	PhoneNumber string `json:"phone" validate:"required,len=10"`
	Email       string `json:"email" validate:"email,required"`
	DateOfBirth string `json:"dateofbirth"`
	Gender      string `json:"gender"`
}

type ForgotPassword struct {
	Email string `json:"email"`
}

type ResetPassword struct {
	Email       string `json:"email"`
	Otp         string `json:"otp"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=12"`
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

// Theater
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
	ID           uint       `json:"id"`
	TheaterID    int        `json:"theater_id"`
	ScreenNumber int        `json:"screen_number"`
	SeatCapacity int        `json:"seat_capacity"`
	ScreenTypeID int        `json:"screen_type_id"`
	ScreenType   ScreenType `json:"ScreenType"`
	Theater      Theater    `json:"Theater"`
}

type Showtime struct {
	ID       uint      `json:"id"`
	MovieID  int       `json:"movie_id"`
	ScreenID int       `json:"screen_id"`
	ShowDate time.Time `json:"show_date"`
	ShowTime time.Time `json:"show_time"`
}

type TheaterWithTypeResponse struct {
	ID              int                 `json:"id"`
	Name            string              `json:"name"`
	Place           string              `json:"place"`
	City            string              `json:"city"`
	District        string              `json:"district"`
	State           string              `json:"state"`
	OwnerID         int                 `json:"owner_id"`
	NumberOfScreens int                 `json:"number_of_screens"`
	TheaterType     TheaterTypeResponse `json:"TheaterType"`
}

type TheaterTypeResponse struct {
	ID              int    `json:"id"`
	TheaterTypeName string `json:"theater_type_name"`
}

type TheatersAndMovieScheduleResponse struct {
	ID         int      `json:"id"`
	MovieID    int      `json:"movie_id"`
	TheaterID  int      `json:"theater_id"`
	ShowtimeID int      `json:"showtime_id"`
	Movie      Movie    `json:"Movie"`
	Theater    Theater  `json:"Theater"`
	Showtime   Showtime `json:"Showtime"`
}

type TheaterResponse struct {
	ID              int                 `json:"id"`
	Name            string              `json:"name"`
	Place           string              `json:"place"`
	City            string              `json:"city"`
	District        string              `json:"district"`
	State           string              `json:"state"`
	NumberOfScreens int                 `json:"number_of_screens"`
	TheaterType     TheaterTypeResponse `json:"TheaterType"`
	MovieSchedules  []MovieSchedule     `json:"MovieSchedule"`
	TheaterScreens  []TheaterScreen     `json:"TheaterScreen"`
}

type MovieSchedule struct {
	ID         int      `json:"id"`
	MovieID    int      `json:"movie_id"`
	TheaterID  int      `json:"theater_id"`
	ShowtimeID int      `json:"showtime_id"`
	Movie      Movie    `json:"Movie"`
	Theater    Theater  `json:"Theater"`
	Showtime   Showtime `json:"Showtime"`
}

type ShowtimeResponse struct {
	ID               uint             `json:"id"`
	MovieID          int              `json:"movie_id"`
	ScreenID         int              `json:"screen_id"`
	ShowDate         time.Time        `json:"show_date"`
	ShowTime         time.Time        `json:"show_time"`
	Movie            Movie            `json:"Movie"`
	TheaterScreenRes TheaterScreenRes `json:"TheaterScreen"`
}

type ListShowTimeResponse struct {
	Theater          Theater            `json:"Theater"`
	TheaterType      TheaterType        `json:"theater_type"`
	ShowtimeResponse []ShowtimeResponse `json:"ShowtimeResponse"`
}

type ListShowTimeByTheaterAndMovie struct {
	Movie                        Movie                          `json:"Movie"`
	Theater                      Theater                        `json:"Theater"`
	ShowtimeResponseWithoutMovie []ShowtimeResponseWithoutMovie `json:"ShowtimeResponse"`
}

type TheaterScreenRes struct {
	ID           uint `json:"id"`
	TheaterID    int  `json:"theater_id"`
	ScreenNumber int  `json:"screen_number"`
	SeatCapacity int  `json:"seat_capacity"`
	ScreenTypeID int  `json:"screen_type_id"`
}

type ShowtimeResponseWithoutMovie struct {
	ID               uint             `json:"id"`
	MovieID          int              `json:"movie_id"`
	ScreenID         int              `json:"screen_id"`
	ShowDate         time.Time        `json:"show_date"`
	ShowTime         time.Time        `json:"show_time"`
	TheaterScreenRes TheaterScreenRes `json:"TheaterScreenRes"`
}

type SeatsByScreenIDRes struct {
	ID                int              `json:"id"`
	ScreenID          int              `json:"screen_id"`
	SeatNumber        string           `json:"seat_number"`
	Row               string           `json:"row"`
	Column            int              `json:"column"`
	SeatCategoryID    int              `json:"seat_category_id"`
	SeatCategoryPrice float64          `json:"seat_category_price"`
	TheaterScreenRes  TheaterScreenRes `json:"TheaterScreen"`
	SeatCategory      SeatCategory     `gorm:"foreignKey:SeatCategoryID"`
}

type ListShowtimesByDateRes struct {
	Theater  TheaterWithTypeResponse `json:"theater_with_type_response"`
	Showtime ShowtimeResponse        `json:"show_time_response"`
}

// Booking
type Booking struct {
	BookingID     uint          `json:"booking_id"`
	UserID        uint          `json:"user_id"`
	ShowtimeID    uint          `json:"showtime_id"`
	BookingDate   time.Time     `json:"booking_date"`
	TotalAmount   float64       `json:"total_amount"`
	PaymentStatus string        `json:"payment_status"`
	BookingSeats  []BookingSeat `json:"booking_seats"`
}

type BookingSeat struct {
	BookingID uint `json:"booking_id"`
	SeatID    uint `json:"seat_id"`
}

type CreateBookingRequest struct {
	UserID      int      `json:"user_id"`
	ShowtimeID  int      `json:"showtime_id"`
	ScreenID    uint     `json:"screen_id"`
	SeatIDs     []uint32 `json:"seat_ids"`
	TotalAmount float64  `json:"total_amount"`
}

type PaymentMethod struct {
	PaymentMethodID uint   `json:"payment_method_id"`
	MethodName      string `json:"method_name"`
}
type Transaction struct {
	TransactionID   uint    `json:"transaction_id"`
	BookingID       uint    `json:"booking_id"`
	UserID          uint    `json:"user_id"`
	PaymentMethodID uint    `json:"payment_method_id"`
	TransactionDate string  `json:"transaction_date"`
	Amount          float64 `json:"amount"`
	OrderID         string  `json:"order_id"`
	Status          string  `json:"status"`
}

type PaymentRequest struct {
	BookingID       uint      `json:"booking_id"`
	UserID          uint      `json:"user_id"`
	PaymentMethodID uint      `json:"payment_method_id"`
	TransactionDate time.Time `json:"transaction_date"`
	Amount          float64   `json:"amount"`
}

type TransactionResponse struct {
	TransactionID   uint    `json:"transaction_id"`
	PaymentMethodID uint    `json:"payment_method_id"`
	TransactionDate string  `json:"transaction_date"`
	Amount          float64 `json:"amount"`
	Status          string  `json:"status"`
}

type PaymentStatusRequest struct {
	OrderID           string `json:"order_id"`
	RazorpayPaymentID string `json:"payment_id"`
}
