package validations

import db "github.com/HouseCham/VetMate/database/sql"

func ValidateUser(user db.InsertNewUserParams) (bool, error) {
	// validation of fullname, password, email and phone fields with generic function
	if isUserValid, err := validateUserOrVet(user.Nombre, user.ApellidoP, user.ApellidoM, user.PasswordHash, user.Email, user.Telefono.String); !isUserValid {
		return false, err
	}
	return true, nil
}