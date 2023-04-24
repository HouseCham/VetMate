package validations

import (
	"errors"
	"regexp"
)


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

	// Check if fullname fields are at least 2 characters long and no more than specified
	if len(nombre) < 2 || len(nombre) > 100 {
		return false, errors.New("nombre must be at least 2 characters long and no more than 100 characters long")
	} else if len(apellidoP) < 2 || len(apellidoP) > 50 {
		return false, errors.New("apellidoP must be at least 2 characters long and no more than 50 characters long")
	} else if len(apellidoM) < 2 || len(apellidoM) > 50 {
		return false, errors.New("apellidoM must be at least 2 characters long and no more than 50 characters long")
	}

	return true, nil
}
// IsValidPasswordInput is a function that checks if a password is valid
// length is between 5 and 72 characters
func isPasswordInputValid(password string) (bool, error) {
	if password == "" {
		return false, errors.New("password is required")
	}
	// Check if password is at least 8 characters long
	if len(password) < 5 || len(password) > 72 {
		return false, errors.New("password length must be between 5 and 72 characters long")
	}
	return true, nil
}