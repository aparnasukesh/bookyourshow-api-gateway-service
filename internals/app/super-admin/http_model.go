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
