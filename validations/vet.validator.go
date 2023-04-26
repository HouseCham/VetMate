package validations

import (
	"errors"

	"github.com/HouseCham/VetMate/interfaces"
)

// ValidateRequest is a function that validates the
// request body for the method specified
// 1 = NewInsert | 2 = Update | 3 = Login
func ValidateRequest(request interfaces.IDatabaseValidation, method int) error {

	switch method {
	// 1 = InsertNewVetOrUser
	case 1:
		// validation of fullname, password, email and phone fields with generic function
		if err := request.ValidateNewRegister(); err != nil {
			return err
		}
		return nil
	// 2 = UpdateVetOrUser
	case 2:
		if err := request.ValidateUpdate(); err != nil {
			return err
		}
		return nil
	// 3 = LoginVetOrUser
	case 3:
		if err := request.ValidateLogin(); err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid method")
	}
}
