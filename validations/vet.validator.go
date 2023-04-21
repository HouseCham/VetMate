package validations

import (
	"errors"

	"github.com/HouseCham/VetMate/database/sql"
	"regexp"
)

// ValidateVet is a function that validates the
// request body for the InsertNewVet function
func ValidateVet(vet db.InsertNewVetParams) (bool, error) {

	// Check if required fields are empty
	if vet.Nombre == "" {
		return false, errors.New("nombre is required")
	} else if vet.ApellidoP == "" {
		return false, errors.New("apellidoP is required")
	}else if vet.ApellidoM == "" {
		return false, errors.New("apellidoM is required")
	} else if vet.PasswordHash == "" {
		return false, errors.New("passwordHash is required")
	}

	// Check fullname fields for special characters
	if hasSpecialCharacters(vet.Nombre) {
		return false, errors.New("nombre cannot contain special characters")
	} else if hasSpecialCharacters(vet.ApellidoP) {
		return false, errors.New("apellidoP cannot contain special characters")
	} else if hasSpecialCharacters(vet.ApellidoM) {
		return false, errors.New("apellidoM cannot contain special characters")
	}

	// Check if fullname fields are at least 2 characters long and no more than specified
	if len(vet.Nombre) < 2 || len(vet.Nombre) > 100 {
		return false, errors.New("nombre must be at least 2 characters long and no more than 100 characters long")
	} else if len(vet.ApellidoP) < 2 || len(vet.ApellidoP) > 50 {
		return false, errors.New("apellidoP must be at least 2 characters long and no more than 50 characters long")
	} else if len(vet.ApellidoM) < 2 || len(vet.ApellidoM) > 50 {
		return false, errors.New("apellidoM must be at least 2 characters long and no more than 50 characters long")
	}

	// Check if password is at least 8 characters long
	if len(vet.PasswordHash) < 8 {
		return false, errors.New("passwordHash must be at least 8 characters long")
	}

	//! TODO: Add validator for email
	if !isValidEmail(vet.Email) {
		return false, errors.New("email is not valid")
	}

	return true, nil
}

// hasSpecialCharacters is a function that checks if a string
// contains special characters
func hasSpecialCharacters(s string) bool {
    re := regexp.MustCompile("[^a-zA-Z0-9]+")
    return re.MatchString(s)
}

// isValidEmail is a function that checks if an email is valid
func isValidEmail(email string) bool {
    // regular expression pattern for email validation
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    // compile the pattern into a regular expression object
    regex := regexp.MustCompile(pattern)
    // match the email against the regular expression
    return regex.MatchString(email)
}
