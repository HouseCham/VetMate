package validations

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/HouseCham/VetMate/config"
)

var Config *config.Config

// ShareConfigFile is a function that shares the
// configuration setted up in main.go
func ShareConfigFile(config *config.Config) {
	Config = config
}


func validateUserOrVet(nombre string, apellidoP string, apellidoM string, password string, email string, telefono string) (bool, error) {
	// Check if fullname is valid
	if isObjectValid, err := isFullNameValid(nombre, apellidoP, apellidoM); !isObjectValid {
		return false, err
	}
	// Check if password is valid
	if isPasswordValid, err := isPasswordInputValid(password); !isPasswordValid {
		return false, err
	}
	// Check if email is valid
	if isEmailValid, err := isEmailValid(email); !isEmailValid {
		return false, err
	}
	// Check if optional fields are not longer than specified	
	if len(telefono) > 20 {
		return false, errors.New("telefono must be no more than 20 characters long")
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
func isFullNameValid(nombre string, apellidoP string, apellidoM string) (bool, error) {
	// Check if required fields are empty
	if nombre == "" {
		return false, errors.New("nombre is required")
	} else if apellidoP == "" {
		return false, errors.New("apellidoP is required")
	}else if apellidoM == "" {
		return false, errors.New("apellidoM is required")
	}
	// Check fullname fields for special characters
	if HasSpecialCharacters(nombre) {
		return false, errors.New("nombre cannot contain special characters")
	} else if HasSpecialCharacters(apellidoP) {
		return false, errors.New("apellidoP cannot contain special characters")
	} else if HasSpecialCharacters(apellidoM) {
		return false, errors.New("apellidoM cannot contain special characters")
	}

	// Check if fullname fields are valid according
	// to the length stablished on config file
	if len(nombre) < Config.DevConfiguration.Parameters.NameMinLength || len(nombre) > Config.DevConfiguration.Parameters.NameMaxLength {
		return false, fmt.Errorf("nombre must be at least %d characters long and no more than %d characters long", Config.DevConfiguration.Parameters.NameMinLength, Config.DevConfiguration.Parameters.NameMaxLength)
	} else if len(apellidoP) < Config.DevConfiguration.Parameters.ApellidoPMinLength || len(apellidoP) > Config.DevConfiguration.Parameters.ApellidoPMaxLength {
		return false, fmt.Errorf("apellidoP must be at least %d characters long and no more than %d characters long", Config.DevConfiguration.Parameters.ApellidoPMinLength, Config.DevConfiguration.Parameters.ApellidoPMaxLength)
	} else if len(apellidoM) < Config.DevConfiguration.Parameters.ApellidoMMinLength || len(apellidoM) > Config.DevConfiguration.Parameters.ApellidoMMaxLength {
		return false, fmt.Errorf("apellidoM must be at least %d characters long and no more than %d characters long", Config.DevConfiguration.Parameters.ApellidoMMinLength, Config.DevConfiguration.Parameters.ApellidoMMaxLength)
	}

	return true, nil
}
// IsValidPasswordInput is a function that checks if a password is valid
// length is between configuration
func isPasswordInputValid(password string) (bool, error) {
	if password == "" {
		return false, errors.New("password is required")
	}
	// Check if password length is valid according to config file
	if len(password) < Config.DevConfiguration.Parameters.PwdMinLength || len(password) > Config.DevConfiguration.Parameters.PwdMaxLength {
		return false, fmt.Errorf("password length must be between %d and %d characters long", Config.DevConfiguration.Parameters.PwdMinLength, Config.DevConfiguration.Parameters.PwdMaxLength)
	}
	return true, nil
}