package controllers

import (
	"database/sql"
	"errors"

	"github.com/HouseCham/VetMate/config"
	db "github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/interfaces"
	"github.com/gofiber/fiber/v2"
)

var DB *sql.DB
var Queries *db.Queries
var Config *config.Config

// ShareDbConnection is a function that shares the
// database connection to all controllers
// so that they can use the same connection
func ShareDbConnection(db *sql.DB) {
	DB = db
	Queries = createNewQuery()
}

// ShareConfigFile is a function that shares the
// configuration setted up in main.go
func ShareConfigFile(config *config.Config) {
	Config = config
}

// createNewQuery is a function that creates a new
// query object for the database connection
func createNewQuery() *db.Queries {
	return db.New(DB)
}

// TrimInputFields is a function that trims all the
// input fields from the request body
// it is an interface function that is used by all
// the controllers
func trimInputFields(input interfaces.INewInsertParams) {
	input.Trim()
}

type IsEmailUsedChan struct {
	Status  int
	Message string
	Err     error
}

// CheckVetEmailAlreadyInUse is a function that checks
// if the email is already in use by checking the database.
// if isUserTable is true, then it means we are trying to check if an user email already exists in database
// otherwise, if it is false, we are trying to check for a vet's email.
func checkEmailAlreadyInUse(isUsedChan chan IsEmailUsedChan, email string, isUserTable bool, c *fiber.Ctx) {
	var emailExists int64
	var err error

	if isUserTable {
		emailExists, err = Queries.CheckUserEmailExists(c.Context(), email)
	} else {
		emailExists, err = Queries.CheckVetEmailExists(c.Context(), email)
	}

	if err != nil {
		isUsedChan <- IsEmailUsedChan{
			Status:  fiber.StatusInternalServerError,
			Message: "Error checking email",
			Err:     err,
		}
	} else if emailExists > 0 {
		isUsedChan <- IsEmailUsedChan{
			Status:  fiber.StatusConflict,
			Message: "Conflict",
			Err:     errors.New("email already in use"),
		}
	} else {
		isUsedChan <- IsEmailUsedChan{
			Status:  fiber.StatusOK,
			Message: "Email not in use",
			Err:     nil,
		}
	}
	close(isUsedChan)
}
