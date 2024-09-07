package user

import (
	"strings"

	"github.com/aparnasukesh/inter-communication/user_admin"
)

func BuildGetUserProfile(res *user_admin.GetProfileResponse) (*UserProfileDetails, error) {
	return &UserProfileDetails{
		Username:    res.ProfileDetails.Username,
		FirstName:   res.ProfileDetails.FirstName,
		LastName:    res.ProfileDetails.LastName,
		PhoneNumber: res.ProfileDetails.Phone,
		DateOfBirth: res.ProfileDetails.DateOfBirth,
		Gender:      res.ProfileDetails.Gender,
	}, nil
}

func ExtractErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	errMsg := err.Error()
	if index := strings.Index(errMsg, "desc = "); index != -1 {
		return errMsg[index+len("desc = "):]
	}
	return errMsg
}
