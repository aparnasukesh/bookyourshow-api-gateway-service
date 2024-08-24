package superadmin

import "time"

type Admin struct {
	Username    string    `json:"username" validate:"required,min=8,max=24"`
	Password    string    `json:"password" validate:"required,min=6,max=12"`
	PhoneNumber string    `json:"phone" validate:"required,len=10"`
	Email       string    `json:"email" validate:"email,required"`
	FirstName   string    `gorm:"not null" json:"firstname" validate:"required,min=4,max=10"`
	LastName    string    `gorm:"not null" json:"lastname" validate:"required,min=4,max=10"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	Otp         string    `json:"otp"`
}

type AdminRequestResponse struct {
	Email string `json:"email" validate:"email,required"`
}

type AdminApproval struct {
	Email      string `json:"email" validate:"email,required"`
	IsVerified bool   `json:"is_verified"`
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
