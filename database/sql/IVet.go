package db

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/HouseCham/VetMate/config"
	"github.com/HouseCham/VetMate/validations"
)

var Config *config.Config

// ShareConfigFile is a function that shares the
// configuration setted up in main.go
func ShareConfigFile(config *config.Config) {
	Config = config
}

func (vet *Veterinario) Trim() {
	vet.Nombre = strings.TrimSpace(vet.Nombre)
	vet.ApellidoP = strings.TrimSpace(vet.ApellidoP)
	vet.ApellidoM = strings.TrimSpace(vet.ApellidoM)
	vet.Email = strings.TrimSpace(vet.Email)
	vet.Telefono.String = strings.TrimSpace(vet.Telefono.String)
	vet.ImgUrl.String = strings.TrimSpace(vet.ImgUrl.String)
	vet.Password = strings.TrimSpace(vet.Password)
}

func (vet *Veterinario) DeleteBlankFields() {
	vet.Nombre = strings.ReplaceAll(vet.Nombre, " ", "")
	vet.ApellidoP = strings.ReplaceAll(vet.ApellidoP, " ", "")
	vet.ApellidoM = strings.ReplaceAll(vet.ApellidoM, " ", "")
	vet.Email = strings.ReplaceAll(vet.Email, " ", "")
	vet.Telefono.String = strings.ReplaceAll(vet.Telefono.String, " ", "")
	vet.ImgUrl.String = strings.ReplaceAll(vet.ImgUrl.String, " ", "")
	vet.Password = strings.ReplaceAll(vet.Password, " ", "")
}

// ValidateUser is a function that validates the
// request body for the InsertNewUser or InsertNewVet function
func (newVetRegister *Veterinario) ValidateNewRegister() error {
	// Check if fullname is valid
	if err := isFullNameValid(newVetRegister.Nombre, newVetRegister.ApellidoP, newVetRegister.ApellidoM); err != nil {
		return err
	}
	// Check if password is valid
	if err := isPasswordInputValid(newVetRegister.Password, Config.DevConfiguration.Parameters.PwdMinLength, Config.DevConfiguration.Parameters.PwdMaxLength); err != nil {
		return err
	}
	// Check if email is valid
	if err := isEmailValid(newVetRegister.Email); err != nil {
		return err
	}
	// Check if optional fields are not longer than specified
	if newVetRegister.Telefono.Valid {
		if err := isPhoneValid(newVetRegister.Telefono.String); err != nil {
			return err
		}
	}
	return nil
}

// ValidateUserUpdate is a function that validates the
// request body for the UpdateUser or UpdateVet function
func (vetUpdate *Veterinario) ValidateUpdate() error {
	if err := isFullNameValid(vetUpdate.Nombre, vetUpdate.ApellidoP, vetUpdate.ApellidoM); err != nil {
		return err
	}
	if vetUpdate.Telefono.Valid {
		if err := isPhoneValid(vetUpdate.Telefono.String); err != nil {
			return err
		}
	}
	if err := isValidImageName(vetUpdate.ImgUrl.String); !err {
		return fmt.Errorf(validations.ErrorMessages["imagen"], "imagenUrl")
	}
	return nil
}

// ValidateUserLogin is a function that validates the
// request body for the VetLogin function
func (vetLogin *Veterinario) ValidateLogin() error {
	if err := isEmailValid(vetLogin.Email); err != nil {
		return err
	} else if err := isPasswordInputValid(vetLogin.Password,
		Config.DevConfiguration.Parameters.PwdMinLength,
		Config.DevConfiguration.Parameters.PwdMaxLength); err != nil {
		return err
	}
	return nil
}

//? COMMON FUNCTIONS

// HasSpecialCharacters is a function that checks if a string
// contains special characters
func hasSpecialCharacters(s string) bool {
	re := regexp.MustCompile("[^a-zA-Z ]+")
	return re.MatchString(s)
}

// IsValidEmail is a function that checks if an email is valid
func isEmailValid(email string) error {
	// regular expression pattern for email validation
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	// compile the pattern into a regular expression object
	regex := regexp.MustCompile(pattern)
	// match the email against the regular expression
	// also check if email is not longer than 150 characters
	if !regex.MatchString(email) {
		return fmt.Errorf(validations.ErrorMessages["correo"], "email")
	} else if len(email) > 150 {
		return fmt.Errorf(validations.ErrorMessages["maximo"], "email", 150)
	}
	return nil
}

// IsFullNameValid is a function that checks if a fullname is valid
// A fullname is valid if it is not empty, does not contain special characters,
// and is at least 2 characters long and no longer than db specified
func isFullNameValid(name string, lastnameP string, lastnameM string) error {
	// Check if required fields are empty
	if name == "" {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "nombre")
	} else if lastnameP == "" {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "apellido paterno")
	} else if lastnameM == "" {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "apellido materno")
	}
	// Check fullname fields for special characters
	if hasSpecialCharacters(name) {
		return fmt.Errorf(validations.ErrorMessages["alfabetico"], "nombre")
	} else if hasSpecialCharacters(lastnameP) {
		return fmt.Errorf(validations.ErrorMessages["alfabetico"], "apellido paterno")
	} else if hasSpecialCharacters(lastnameM) {
		return fmt.Errorf(validations.ErrorMessages["alfabetico"], "apellido materno")
	}
	// Check if fullname fields are valid according
	// to the length stablished on config file
	if len(name) < Config.DevConfiguration.Parameters.NameMinLength || len(name) > Config.DevConfiguration.Parameters.NameMaxLength {
		return fmt.Errorf(validations.ErrorMessages["min_max"], "nombre", Config.DevConfiguration.Parameters.NameMinLength, Config.DevConfiguration.Parameters.NameMaxLength)
	} else if len(lastnameP) < Config.DevConfiguration.Parameters.ApellidoPMinLength || len(lastnameP) > Config.DevConfiguration.Parameters.ApellidoPMaxLength {
		return fmt.Errorf(validations.ErrorMessages["min_max"], "apellido paterno", Config.DevConfiguration.Parameters.ApellidoPMinLength, Config.DevConfiguration.Parameters.ApellidoPMaxLength)
	} else if len(lastnameM) < Config.DevConfiguration.Parameters.ApellidoMMinLength || len(lastnameM) > Config.DevConfiguration.Parameters.ApellidoMMaxLength {
		return fmt.Errorf(validations.ErrorMessages["min_max"], "apellido materno", Config.DevConfiguration.Parameters.ApellidoMMinLength, Config.DevConfiguration.Parameters.ApellidoMMaxLength)
	}
	return nil
}

// IsValidPasswordInput is a function that checks if a password is valid
// length is between configuration
func isPasswordInputValid(password string, pwdMinLength int, pwdMaxLength int) error {
	if password == "" {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "contraseña")
	}
	// Check if password length is valid according to config file
	if len(password) < pwdMinLength || len(password) > pwdMaxLength {
		return fmt.Errorf(validations.ErrorMessages["min_max"], "contraseña", pwdMinLength, pwdMaxLength)
	}
	return nil
}

// IsPhoneValid is a function that checks if a phone number is valid
// A phone number is valid if it is not longer than 20 characters
// and not shorter than 5 characters
func isPhoneValid(phone string) error {
	if len(phone) < 5 || len(phone) > 20 {
		return fmt.Errorf(validations.ErrorMessages["min_max"], "teléfono", 5, 20)
	} else if err := isValidInt(phone, "teléfono"); err != nil {
		return err
	}
	return nil
}

// IsValidImageName is a function that checks if an image name is valid
// A valid image name is a string that contains only letters, numbers,
// dashes and underscores, and ends with a valid image extension
func isValidImageName(imageName string) bool {
	pattern := "^[a-zA-Z0-9-_]+\\.[a-zA-Z]{2,4}$"
	match, _ := regexp.MatchString(pattern, imageName)
	return match
}

// isValidAndNoNegativeInt is a function that checks if a string
// is a valid integer and if it is not negative
func isValidInt(numberStr string, fieldName string) error {
	re := regexp.MustCompile("^[0-9]+$")
	if number, err := strconv.Atoi(numberStr); err != nil && number < 0 {
		return fmt.Errorf(validations.ErrorMessages["numerico"], fieldName)
	}
	if !re.MatchString(numberStr) {
		return fmt.Errorf(validations.ErrorMessages["numerico"], fieldName)
	}
	return nil
}

// hasOnlyAlphanumericAndPunctuation is a function that checks if a string
// contains only letters, numbers and punctuation signs
func hasOnlyAlphanumericAndPunctuation(str string) bool {
	// Regular expression pattern for letters and punctuation signs
	pattern := "^[a-zA-Z0-9[:punct:]]+$"
	// Compile the regular expression
	regex := regexp.MustCompile(pattern)
	// Test if the string matches the pattern
	isValid := regex.MatchString(str)
	return isValid
}

// isValidMySQLDate is a function that checks if a string is a valid MySQL date
func isValidMySQLDate(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// isValidWeightValue is a function that checks if a string is a valid weight value, with 2 decimals
func isValidWeightValue(weightStr string) bool {
	regex := `^\d+(\.\d{2})?$`
	match, _ := regexp.MatchString(regex, weightStr)
	return match
}

// parseWeightValue is a function that parses a string to a float64
// value
func parseWeightValue(weightStr string) (float64, error) {
	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		return 0, err
	}
	return weight, nil
}