package user

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

func ValidateUser(user User) error {
	validate := validator.New()

	err := validate.Struct(user)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make([]string, len(validationErrors))

		for i, validationErr := range validationErrors {
			fieldName := validationErr.Field()
			switch fieldName {
			case "Email":
				errorMessages[i] = "Invalid Email"
			case "Username":
				errorMessages[i] = "Invalid Username, Minimum 8 letters or Maximum 24 letters required"
			case "FirstName":
				errorMessages[i] = "Invalid Firstname, Minimum 4 letters or Maximum 10 letters required"
			case "LastName":
				errorMessages[i] = "Invalid Lastname, Minimum 4 letters or Maximum 10 letters required"
			case "Password":
				errorMessages[i] = "Invalid password, Minimum 6 letters or Maximum 12 letters required"
			case "PhoneNumber":
				errorMessages[i] = "Invalid Phone Number"
			default:
				errorMessages[i] = "Validation failed"
			}
		}

		return fmt.Errorf(strings.Join(errorMessages, ", "))
	}
	return nil
}
