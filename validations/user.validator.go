package validations

import db "github.com/HouseCham/VetMate/database/sql"

func ValidateUser(user db.InsertNewUserParams) (bool, error) {
	// validation of fullname, password, email and phone fields with generic function
	
	return true, nil
}
