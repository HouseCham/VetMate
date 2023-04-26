package db

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/HouseCham/VetMate/config"
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
	vet.PasswordHash = strings.TrimSpace(vet.PasswordHash)
}

// ValidateUser is a function that validates the
// request body for the InsertNewUser or InsertNewVet function
func (newVetRegister *Veterinario) ValidateNewRegister() error {
	// Check if fullname is valid
	if err := isFullNameValid(newVetRegister.Nombre, newVetRegister.ApellidoP, newVetRegister.ApellidoM); err != nil {
		return err
	}
	// Check if password is valid
	if err := isPasswordInputValid(newVetRegister.PasswordHash, Config.DevConfiguration.Parameters.PwdMinLength, Config.DevConfiguration.Parameters.PwdMaxLength); err != nil {
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
		return errors.New("campo imgUrl no válido")
	}
	return nil
}

// ValidateUserLogin is a function that validates the
// request body for the VetLogin function
func (vetLogin *Veterinario) ValidateLogin() error {
	if err := isEmailValid(vetLogin.Email); err != nil {
		return err
	} else if err := isPasswordInputValid(vetLogin.PasswordHash,
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
		return errors.New("campo email no válido")
	} else if len(email) > 150 {
		return errors.New("campo email no debe superar los 150 caracteres")
	}
	return nil
}

// IsFullNameValid is a function that checks if a fullname is valid
// A fullname is valid if it is not empty, does not contain special characters,
// and is at least 2 characters long and no longer than db specified
func isFullNameValid(name string, lastnameP string, lastnameM string) error {
	// Check if required fields are empty
	if name == "" {
		return errors.New("campo nombre es requerido")
	} else if lastnameP == "" {
		return errors.New("campo apellido paterno es requerido")
	} else if lastnameM == "" {
		return errors.New("campo apellido materno es requerido")
	}
	// Check fullname fields for special characters
	if hasSpecialCharacters(name) {
		return errors.New("campo nombre no puede contener caracteres especiales")
	} else if hasSpecialCharacters(lastnameP) {
		return errors.New("campo apellido paterno no puede contener caracteres especiales")
	} else if hasSpecialCharacters(lastnameM) {
		return errors.New("campo apellido materno no puede contener caracteres especiales")
	}
	// Check if fullname fields are valid according
	// to the length stablished on config file
	if len(name) < Config.DevConfiguration.Parameters.NameMinLength || len(name) > Config.DevConfiguration.Parameters.NameMaxLength {
		return fmt.Errorf("campo nombre debe tener entre %d y %d caracteres", Config.DevConfiguration.Parameters.NameMinLength, Config.DevConfiguration.Parameters.NameMaxLength)
	} else if len(lastnameP) < Config.DevConfiguration.Parameters.ApellidoPMinLength || len(lastnameP) > Config.DevConfiguration.Parameters.ApellidoPMaxLength {
		return fmt.Errorf("campo apellido paterno debe tener entre %d y %d caracteres", Config.DevConfiguration.Parameters.ApellidoPMinLength, Config.DevConfiguration.Parameters.ApellidoPMaxLength)
	} else if len(lastnameM) < Config.DevConfiguration.Parameters.ApellidoMMinLength || len(lastnameM) > Config.DevConfiguration.Parameters.ApellidoMMaxLength {
		return fmt.Errorf("campo apellido materno debe tener entre %d y %d caracteres", Config.DevConfiguration.Parameters.ApellidoMMinLength, Config.DevConfiguration.Parameters.ApellidoMMaxLength)
	}
	return nil
}

// IsValidPasswordInput is a function that checks if a password is valid
// length is between configuration
func isPasswordInputValid(password string, pwdMinLength int, pwdMaxLength int) error {
	if password == "" {
		return errors.New("campo contraseña es requerido")
	}
	// Check if password length is valid according to config file
	if len(password) < pwdMinLength || len(password) > pwdMaxLength {
		return fmt.Errorf("campo contraseña debe tener entre %d y %d caracteres", pwdMinLength, pwdMaxLength)
	}
	return nil
}

// IsPhoneValid is a function that checks if a phone number is valid
// A phone number is valid if it is not longer than 20 characters
// and not shorter than 5 characters
func isPhoneValid(phone string) error {
	if len(phone) < 5 || len(phone) > 20 {
		return errors.New("campo teléfono debe tener entre 5 y 20 caracteres")
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
		return fmt.Errorf("campo %s inválido", fieldName)
	}
	if !re.MatchString(numberStr) {
		return fmt.Errorf("campo %s inválido", fieldName)
	}
	return nil
}
