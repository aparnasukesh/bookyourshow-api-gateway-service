package admin

type Admin struct {
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

type AdminProfileDetails struct {
	Username    string `json:"username" validate:"required,min=8,max=24"`
	FirstName   string `gorm:"not null" json:"firstname" validate:"required,min=4,max=10"`
	LastName    string `gorm:"not null" json:"lastname" validate:"required,min=4,max=10"`
	Phone       string `json:"phone" validate:"required,len=10"`
	Email       string `json:"email" validate:"email,required"`
	Dateofbirth string `json:"dateofbirth"`
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
	ID                int     `json:"id"`
	SeatCategoryName  string  `json:"seat_category_name"`
	SeatCategoryPrice float64 `json:"seat_category_price"`
}

type Theater struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Location        string `json:"location"`
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
