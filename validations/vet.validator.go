package validations

import (
	"errors"

	db "github.com/HouseCham/VetMate/database/sql"
)

// ValidateVet is a function that validates the
// request body for the InsertNewVet function
func ValidateVet(vet db.InsertNewVetParams) (bool, error) {
	// validation of fullname, password, email and phone fields with generic function
	if isVetValid, err := validateUserOrVet(vet.Nombre, vet.ApellidoP, vet.ApellidoM, vet.PasswordHash, vet.Email, vet.Telefono.String, 2, 72); !isVetValid {
		return false, err
	}
	// Check if optional fields are not longer than specified
	if len(vet.ImgUrl.String) > 255 {
		return false, errors.New("img url must be no more than 255 characters long")
	}
	return true, nil
}

// ValidateVetLogin is a function that validates the
// request body for the VetLogin function
// It checks if the email and password are valid
func ValidateVetLogin(email string, password string) (bool, error) {
	if isValid, err := isEmailValid(email); !isValid {
		return false, err
	} else if isValidPwd, err := isPasswordInputValid(password,
		Config.DevConfiguration.Parameters.PwdMinLength,
		Config.DevConfiguration.Parameters.PwdMaxLength); !isValidPwd {
		return false, err
	}
	return true, nil
}
