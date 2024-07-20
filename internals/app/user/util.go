package user

import "github.com/aparnasukesh/inter-communication/user_admin"

func BuildGetUserProfile(res *user_admin.GetProfileResponse) (*UserProfileDetails, error) {
	return &UserProfileDetails{
		Username:    res.ProfileDetails.Username,
		FirstName:   res.ProfileDetails.FirstName,
		LastName:    res.ProfileDetails.LastName,
		Phone:       res.ProfileDetails.Phone,
		Email:       res.ProfileDetails.Email,
		Dateofbirth: res.ProfileDetails.DateOfBirth,
		Gender:      res.ProfileDetails.Gender,
	}, nil
}
