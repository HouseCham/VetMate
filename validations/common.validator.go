package validations

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/HouseCham/VetMate/config"
	db "github.com/HouseCham/VetMate/database/sql"
)

var Config *config.Config

// ShareConfigFile is a function that shares the
// configuration setted up in main.go
func ShareConfigFile(config *config.Config) {
	Config = config
}

// ValidateUser is a function that validates the
// request body for the InsertNewUser or InsertNewVet function
func validateUserOrVetRegister(vet *db.Veterinario) (bool, error) {
	// Check if fullname is valid
	if isFullnameValid, err := isFullNameValid(vet); !isFullnameValid {
		return false, err
	}
	// Check if password is valid
	if isPasswordValid, err := isPasswordInputValid(vet.PasswordHash, Config.DevConfiguration.Parameters.PwdMinLength, Config.DevConfiguration.Parameters.PwdMaxLength); !isPasswordValid {
		return false, err
	}
	// Check if email is valid
	if isEmailValid, err := isEmailValid(vet.Email); !isEmailValid {
		return false, err
	}
	// Check if optional fields are not longer than specified
	if isPhoneValid, err := isPhoneValid(vet); !isPhoneValid {
		return false, err
	}
	return true, nil
}

// ValidateUserUpdate is a function that validates the
// request body for the UpdateUser or UpdateVet function
func validateUserOrVetUpdate(vet *db.Veterinario) (bool, error) {
	if isFullnameValid, err := isFullNameValid(vet); !isFullnameValid {
		return false, err
	}
	if isPhoneValid, err := isPhoneValid(vet); !isPhoneValid {
		return false, err
	}
	if !isValidImageName(vet.ImgUrl.String) {
		return false, errors.New("img url is not valid")
	}
	return true, nil
}

// HasSpecialCharacters is a function that checks if a string
// contains special characters
func HasSpecialCharacters(s string) bool {
	re := regexp.MustCompile("[^a-zA-Z]+")
	return re.MatchString(s)
}

// IsValidEmail is a function that checks if an email is valid
func isEmailValid(email string) (bool, error) {
	// regular expression pattern for email validation
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	// compile the pattern into a regular expression object
	regex := regexp.MustCompile(pattern)
	// match the email against the regular expression
	// also check if email is not longer than 150 characters
	if !regex.MatchString(email) {
		return false, errors.New("email is not valid")
	} else if len(email) > 150 {
		return false, errors.New("email length is greater than 150 characters")
	}
	return true, nil
}

// IsFullNameValid is a function that checks if a fullname is valid
// A fullname is valid if it is not empty, does not contain special characters,
// and is at least 2 characters long and no longer than db specified
func isFullNameValid(vet *db.Veterinario) (bool, error) {
	// Check if required fields are empty
	if vet.Nombre == "" {
		return false, errors.New("nombre is required")
	} else if vet.ApellidoP == "" {
		return false, errors.New("apellidoP is required")
	} else if vet.ApellidoM == "" {
		return false, errors.New("apellidoM is required")
	}
	// Check fullname fields for special characters
	if HasSpecialCharacters(vet.Nombre) {
		return false, errors.New("nombre cannot contain special characters")
	} else if HasSpecialCharacters(vet.ApellidoP) {
		return false, errors.New("apellidoP cannot contain special characters")
	} else if HasSpecialCharacters(vet.ApellidoM) {
		return false, errors.New("apellidoM cannot contain special characters")
	}
	// Check if fullname fields are valid according
	// to the length stablished on config file
	if len(vet.Nombre) < Config.DevConfiguration.Parameters.NameMinLength || len(vet.Nombre) > Config.DevConfiguration.Parameters.NameMaxLength {
		return false, fmt.Errorf("nombre must be at least %d characters long and no more than %d characters long", Config.DevConfiguration.Parameters.NameMinLength, Config.DevConfiguration.Parameters.NameMaxLength)
	} else if len(vet.ApellidoP) < Config.DevConfiguration.Parameters.ApellidoPMinLength || len(vet.ApellidoP) > Config.DevConfiguration.Parameters.ApellidoPMaxLength {
		return false, fmt.Errorf("apellidoP must be at least %d characters long and no more than %d characters long", Config.DevConfiguration.Parameters.ApellidoPMinLength, Config.DevConfiguration.Parameters.ApellidoPMaxLength)
	} else if len(vet.ApellidoM) < Config.DevConfiguration.Parameters.ApellidoMMinLength || len(vet.ApellidoM) > Config.DevConfiguration.Parameters.ApellidoMMaxLength {
		return false, fmt.Errorf("apellidoM must be at least %d characters long and no more than %d characters long", Config.DevConfiguration.Parameters.ApellidoMMinLength, Config.DevConfiguration.Parameters.ApellidoMMaxLength)
	}
	return true, nil
}

// IsValidPasswordInput is a function that checks if a password is valid
// length is between configuration
func isPasswordInputValid(password string, pwdMinLength int, pwdMaxLength int) (bool, error) {
	if password == "" {
		return false, errors.New("password is required")
	}
	// Check if password length is valid according to config file
	if len(password) < pwdMinLength || len(password) > pwdMaxLength {
		return false, fmt.Errorf("password length must be between %d and %d characters long", pwdMinLength, pwdMaxLength)
	}
	return true, nil
}

// IsPhoneValid is a function that checks if a phone number is valid
// A phone number is valid if it is not longer than 20 characters
// and not shorter than 5 characters
func isPhoneValid(vet *db.Veterinario) (bool, error) {
	if vet.Telefono.Valid {
		if len(vet.Telefono.String) < 5 || len(vet.Telefono.String) > 20 {
			return false, errors.New("telefono must be at least 5 and no more than 20 characters long")
		}
	}
	return true, nil
}

// IsValidImageName is a function that checks if an image name is valid
// A valid image name is a string that contains only letters, numbers,
// dashes and underscores, and ends with a valid image extension
func isValidImageName(imageName string) bool {
	pattern := "^[a-zA-Z0-9-_]+\\.[a-zA-Z]{2,4}$"
	match, _ := regexp.MatchString(pattern, imageName)
	return match
}
