package validations

import (
	"errors"

	db "github.com/HouseCham/VetMate/database/sql"
)

// ValidateVet is a function that validates the
// request body for the InsertNewVet function
// method parameter specifies the method that is calling this function
// 1 = InsertNewVet | 2 = UpdateVet | 3 = DeleteVet
func ValidateVet(vet *db.Veterinario, method int) (bool, error) {

	switch method {
	// 1 = InsertNewVet
	case 1:
		// validation of fullname, password, email and phone fields with generic function
		if isVetValid, err := validateUserOrVetLogin(vet, vet.PasswordHash, vet.Email, vet.Telefono.String, 2, 72); !isVetValid {
			return false, err
		}
		return true, nil
	// 2 = UpdateVet
	case 2:
		if isVetValid, err := validateUserOrVetUpdate(vet); !isVetValid {
			return false, err
		}
		return true, nil
	// 3 = DeleteVet
	case 3:
		break
	default:
		return false, errors.New("invalid method")
	}

	// Check if optional fields are not longer than specified
	if len(vet.ImgUrl.String) > 255 {
		return false, errors.New("img url must be no more than 255 characters long")
	}

	// Check if dates are not set
	if vet.FechaRegistro.Time.String() != "" || vet.FechaDelete.Time.String() != "" || vet.FechaUpdate.Time.String() != "" {
		return false, errors.New("vet cannot set dates")
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
