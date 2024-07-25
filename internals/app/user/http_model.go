package user

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
	Phone       string `json:"phone" validate:"required,len=10"`
	Email       string `json:"email" validate:"email,required"`
	Dateofbirth string `json:"dateofbirth"`
	Gender      string `json:"gender"`
}
