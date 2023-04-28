package db

import (
	"fmt"
	"strings"

	"github.com/HouseCham/VetMate/validations"
)



func (user *Usuario) Trim(){
	user.Nombre = strings.TrimSpace(user.Nombre)
	user.ApellidoP = strings.TrimSpace(user.ApellidoP)
	user.ApellidoM = strings.TrimSpace(user.ApellidoM)
	user.Email = strings.TrimSpace(user.Email)
	user.Telefono.String = strings.TrimSpace(user.Telefono.String)
	user.PasswordHash = strings.TrimSpace(user.PasswordHash)
	user.Calle = strings.TrimSpace(user.Calle)
	user.Colonia = strings.TrimSpace(user.Colonia)
	user.Ciudad = strings.TrimSpace(user.Ciudad)
	user.Estado = strings.TrimSpace(user.Estado)
	user.Cp = strings.TrimSpace(user.Cp)
	user.Pais = strings.TrimSpace(user.Pais)
	user.NumExt = strings.TrimSpace(user.NumExt)
	user.NumInt.String = strings.TrimSpace(user.NumInt.String)
	user.Referencia.String = strings.TrimSpace(user.Referencia.String)
}

// ValidateNewRegister is a function that validates the
// request body for the InsertNewUser or InsertNewVet function
// this function is implemented from IDatabase interface
func (newUserRegister *Usuario) ValidateNewRegister() error {
	// Check if fullname is valid
	if err := isFullNameValid(newUserRegister.Nombre, newUserRegister.ApellidoP, newUserRegister.ApellidoM); err != nil {
		return err
	}
	// Check if password is valid
	if err := isPasswordInputValid(newUserRegister.PasswordHash, Config.DevConfiguration.Parameters.PwdMinLength, Config.DevConfiguration.Parameters.PwdMaxLength); err != nil {
		return err
	}
	// Check if email is valid
	if err := isEmailValid(newUserRegister.Email); err != nil {
		return err
	}
	// Check if optional fields are not longer than specified
	if newUserRegister.Telefono.Valid {
		if err := isPhoneValid(newUserRegister.Telefono.String); err != nil {
			return err
		}
	}
	// Check if address is valid
	if err := isAddressValid(newUserRegister.Calle, newUserRegister.Colonia, newUserRegister.Ciudad, newUserRegister.Estado, newUserRegister.Cp, newUserRegister.NumExt, newUserRegister.NumInt.String); err != nil {
		return err
	}

	// Check if reference is not longer than specified
	if newUserRegister.Referencia.Valid {
		if len(newUserRegister.Referencia.String) > 255 {
			return fmt.Errorf(validations.ErrorMessages["maximo"], "referencia", 255)
		}
	}

	return nil
}

// ValidateUpdate is a function that validates the
// request body for the UpdateUser or UpdateVet function
// this function is implemented from IDatabase interface
func(userUpdate *Usuario) ValidateUpdate() error {
	// Check if fullname is valid
	if err := isFullNameValid(userUpdate.Nombre, userUpdate.ApellidoP, userUpdate.ApellidoM); err != nil {
		return err
	}
	
	// Check if optional fields are not longer than specified
	if userUpdate.Telefono.Valid {
		if err := isPhoneValid(userUpdate.Telefono.String); err != nil {
			return err
		}
	}
	// Check if address is valid
	if err := isAddressValid(userUpdate.Calle, userUpdate.Colonia, userUpdate.Ciudad, userUpdate.Estado, userUpdate.Cp, userUpdate.NumExt, userUpdate.NumInt.String); err != nil {
		return err
	}

	// Check if reference is not longer than specified
	if userUpdate.Referencia.Valid {
		if len(userUpdate.Referencia.String) > 255 {
			return fmt.Errorf(validations.ErrorMessages["maximo"], "referencia", 255)
		}
	}
	return nil
}

// ValidateLogin is a function that validates the
// request body for the LoginUser or LoginVet function
// this function is implemented from IDatabase interface
func(userLogin *Usuario) ValidateLogin() error {
	if err := isEmailValid(userLogin.Email); err != nil {
		return err
	} else if err := isPasswordInputValid(userLogin.PasswordHash,
		Config.DevConfiguration.Parameters.PwdMinLength,
		Config.DevConfiguration.Parameters.PwdMaxLength); err != nil {
		return err
	}
	return nil
}

// isAddressValid is a function that validates the address
// fields for the InsertNewUser function
func isAddressValid(calle string, neighborhood string, city string, state string, zipcode string, extNum string, intNum string) error {
	// Check if not null values are empty
	if calle == "" {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "calle")
	} else if neighborhood == "" {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "colonia")
	} else if city == "" {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "ciudad")
	} else if state == "" {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "estado")
	} else if zipcode == "" {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "código postal")
	} else if extNum == "" {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "número exterior")
	}

	// Check if values are not longer than specified
	if len(calle) > Config.DevConfiguration.Parameters.Address.StreetMaxLength {
		return fmt.Errorf(validations.ErrorMessages["maximo"], "calle", Config.DevConfiguration.Parameters.Address.StreetMaxLength)
	} else if len(neighborhood) > Config.DevConfiguration.Parameters.Address.NeighborhoodMaxLength {
		return fmt.Errorf(validations.ErrorMessages["maximo"], "colonia", Config.DevConfiguration.Parameters.Address.NeighborhoodMaxLength)
	} else if len(city) > Config.DevConfiguration.Parameters.Address.CityMaxLength {
		return fmt.Errorf(validations.ErrorMessages["maximo"], "ciudad", Config.DevConfiguration.Parameters.Address.CityMaxLength)
	} else if len(state) > Config.DevConfiguration.Parameters.Address.StateMaxLength {
		return fmt.Errorf(validations.ErrorMessages["maximo"], "estado", Config.DevConfiguration.Parameters.Address.StateMaxLength)
	} else if len(zipcode) > Config.DevConfiguration.Parameters.Address.ZipCodeMaxLength {
		return fmt.Errorf(validations.ErrorMessages["maximo"], "código postal", Config.DevConfiguration.Parameters.Address.ZipCodeMaxLength)
	} else if len(extNum) > Config.DevConfiguration.Parameters.Address.ExtNumberMaxLength {
		return fmt.Errorf(validations.ErrorMessages["maximo"], "número interior", Config.DevConfiguration.Parameters.Address.ExtNumberMaxLength)
	} else if len(intNum) > Config.DevConfiguration.Parameters.Address.IntNumberMaxLength {
		return fmt.Errorf(validations.ErrorMessages["maximo"], "número exterior", Config.DevConfiguration.Parameters.Address.IntNumberMaxLength)
	}

	// Check if address numbers are valid
	if err := isValidInt(extNum, "número ext"); err != nil {
		return err
    } else if err := isValidInt(zipcode, "código postal"); err != nil {
		return err
	} else if intNum != "" {
		if err := isValidInt(intNum, "número int"); err != nil {
			return err
		}
	}

	// Check if string values have special characters
	if hasSpecialCharacters(calle) {
		return fmt.Errorf(validations.ErrorMessages["alfabetico"], "calle")
	} else if hasSpecialCharacters(neighborhood) {
		return fmt.Errorf(validations.ErrorMessages["alfabetico"], "colonia")
	} else if hasSpecialCharacters(city) {
		return fmt.Errorf(validations.ErrorMessages["alfabetico"], "ciudad")
	} else if hasSpecialCharacters(state) {
		return fmt.Errorf(validations.ErrorMessages["alfabetico"], "estado")
	}

	return nil
}