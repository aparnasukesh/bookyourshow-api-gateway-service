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

type Movie struct {
	Title       string  `gorm:"type:varchar(100);not null"`
	Description string  `gorm:"type:text"`
	Duration    int     `gorm:"not null"`
	Genre       string  `gorm:"type:varchar(50)"`
	ReleaseDate string  `gorm:"not null"`
	Rating      float64 `gorm:"type:decimal(3,1)"`
	Language    string  `gorm:"type:varchar(100);not null"`
}
